package server

import (
	"encoding/base64"
	"errors"
	"main/client"
	"strings"

	"github.com/gin-gonic/gin"
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
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"msg":   token,
		"creds": credentials,
	})
}

func getArtifactSHA(c *gin.Context) {
	var token []string
	host := c.DefaultQuery("host", serverConfig.Host)

	if token = c.Request.Header["Token"]; token == nil {
		c.JSON(400, gin.H{
			"msg": "no Token header",
		})
		return
	}

	aux := c.DefaultQuery("image", "")
	project, image, err := splitImage(aux)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	sha, err := client.GetArtifactSHA(host, "v2.0/", token[0], project, image)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"image":   image,
		"project": project,
		"sha":     sha,
	})
}

func checkArtifactSHA(c *gin.Context) {
	var token []string

	if token = c.Request.Header["Token"]; token == nil {
		c.JSON(400, gin.H{
			"msg": "no Token header",
		})
		return
	}

	aux := c.DefaultQuery("image", "")
	project, image, err := splitImage(aux)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}

	targetDigest := c.DefaultQuery("targetDigest", "")

	host := c.DefaultQuery("host", serverConfig.Host)
	sha, err := client.GetArtifactSHA(host, "v2.0/", token[0], project, image)
	if err != nil {
		c.JSON(400, gin.H{
			"msg": err.Error(),
		})
		return
	}
	equals := strings.EqualFold(targetDigest, sha)
	code := 200
	if !equals {
		code = 400
	}

	c.JSON(code, gin.H{
		"image":        image,
		"project":      project,
		"sha":          sha,
		"targetDigset": targetDigest,
		"equals":       equals,
	})
}

func Execute(config ServerConfig) {
	serverConfig = config
	route := gin.Default()
	route.GET("/jwt", getToken)
	route.GET("/artifact/sha", getArtifactSHA)
	route.GET("/artifact/check_sha", checkArtifactSHA)
	route.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
