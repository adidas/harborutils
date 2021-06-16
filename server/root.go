package server

import (
	"encoding/base64"
	"errors"
	"main/client"
	"net/http"
	"strings"

	_ "main/server/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type ServerConfig struct {
	ClientId string
	TenantId string
	Host     string
}

var serverConfig ServerConfig

func basicAuth(credentials string) (string, string) {
	auth := strings.SplitN(credentials, " ", 2)

	if len(auth) != 2 || auth[0] != "Basic" {
		return "", ""
	}
	payload, _ := base64.StdEncoding.DecodeString(auth[1])
	pair := strings.SplitN(string(payload), ":", 2)

	if len(pair) != 2 {
		return "", ""
	}
	return pair[0], pair[1]
}

func splitImage(image string) (string, string, error) {
	aux := strings.SplitN(image, "/", 2)

	if len(aux) != 2 {
		return "", "", errors.New("No proyect or artifact")
	}
	return aux[0], aux[1], nil

}

type ParamsGetArtifactSha struct {
	Name string `uri:"name" binding:"required"`
}

//
// @Summary Get Bearer to use harborUtils or Harbor Api
// @Description get Bearer, using https://github.com/goharbor/harbor/issues/13683#issuecomment-739036574
// @Accept  json
// @Produce  json
// @Param   client_id     query    string     false        "Oidc client id for authentication"
// @Param   tenant_id     query    string     false        "Azure tenant for oidc authentication"
// @Security BasicAuth
// @Success 200 {object} server.Token	"Success"
// @Failure 400 {object} server.APIError "Bad request"
// @Router /jwt [get]
func getToken(c *gin.Context) {
	username, password := "", ""
	credentials := c.Request.Header["Authorization"]
	if credentials == nil {
		c.JSON(400, gin.H{
			"msg": "no Authorization header",
		})
		return
	}
	username, password = basicAuth(credentials[0])
	// username, password := basicAuth(credentials)
	clientId := c.DefaultQuery("client_id", serverConfig.ClientId)
	tenant := c.DefaultQuery("tenant_id", serverConfig.TenantId)
	token, err := client.GetOidcBearer(clientId, tenant, "", username, password)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, Token{Token: token})
}

//
// @Summary Get image digest from Harbor
// @Description get image digest from Harbor, harbor api: /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}
// @Accept  json
// @Produce  json
// @Param   host     query    string     false        "Harbor url"
// @Param   image     query    string     true        "image name"
// @Param   Token     header    string     true        "Bearer to use harbor api"
// @Success 200 {object} server.ArtifactSha	"Success"
// @Failure 400 {object} server.APIError "Bad request"
// @Router /artifact/sha [get]
func getArtifactSHA(c *gin.Context) {
	var token []string
	host := c.DefaultQuery("host", serverConfig.Host)

	if token = c.Request.Header["Token"]; token == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "no Token header",
		})
		return
	}

	aux := c.DefaultQuery("image", "")
	project, image, err := splitImage(aux)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{Msg: err.Error()})
		return
	}

	sha, err := client.GetArtifactSHA(host, "v2.0/", token[0], project, image)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{Msg: err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, ArtifactSha{Image: image, Project: project, Sha: sha})
}

//
// @Summary Check image digest from Harbor
// @Description Check image digest from Harbor, harbor api: /projects/{project_name}/repositories/{repository_name}/artifacts/{reference}
// @Accept  json
// @Produce  json
// @Param   host     query    string     false        "Harbor url"
// @Param   image     query    string     true        "image name"
// @Param   targetDigest     query    string     true        "sha digest"
// @Param   Token     header    string     true        "Bearer to use harbor api"
// @Success 200 {object} server.ArtifactCheckSha	"Success"
// @Failure 400 {object} server.APIError "Bad request"
// @Router /artifact/check_sha [get]
func checkArtifactSHA(c *gin.Context) {
	var token []string

	if token = c.Request.Header["Token"]; token == nil {
		c.JSON(http.StatusBadRequest, APIError{Msg: "no Token header"})
		return
	}

	aux := c.DefaultQuery("image", "")
	project, image, err := splitImage(aux)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{Msg: err.Error()})
		return
	}

	targetDigest := c.DefaultQuery("targetDigest", "")

	host := c.DefaultQuery("host", serverConfig.Host)
	sha, err := client.GetArtifactSHA(host, "v2.0/", token[0], project, image)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIError{Msg: err.Error()})
		return
	}
	equals := strings.EqualFold(targetDigest, sha)
	code := http.StatusAccepted
	if !equals {
		code = http.StatusBadRequest
	}

	c.JSON(code, gin.H{
		"image":        image,
		"project":      project,
		"sha":          sha,
		"targetDigset": targetDigest,
		"equals":       equals,
	})
}

//
// @Summary Health check API
// @Description The endpoint returns the health stauts of the system.
// @Produce  json
// @Success 200 {object} server.HealthStatus	"Success"
// @Router /health [get]
func health(c *gin.Context) {
	c.JSON(http.StatusOK, HealthStatus{Status: "healthy"})
}

// @title HarborUtils API
// @version 1.0
// @description These APIs provide services for using HarborUtiuls.
// @contact.url https://*****/confluence/spaces/viewspace.action?key=CICDTOOLS
// @securityDefinitions.basic BasicAuth
func Execute(config ServerConfig) {
	serverConfig = config
	route := gin.Default()
	route.GET("/jwt", getToken)
	route.GET("/artifact/sha", getArtifactSHA)
	route.GET("/artifact/check_sha", checkArtifactSHA)
	route.GET("/health", health)

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	route.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
