package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

func getArtifactSHA(clientId, tenant, server, user, password, apiVersion, project, image string) (string, error) {
	artifact := strings.Split(image, ":")
	if len(artifact) < 2 {
		artifact = append(artifact, "latest")
	}
	idtoken, err := getOidcBearer(clientId, tenant, server, user, password)

	if err != nil {
		log.Fatal("There were an error getting idtoken")
	}
	//log.Println(idtoken)
	bearer := fmt.Sprintf("Bearer %v", idtoken)
	url := fmt.Sprintf("%v/api/%vprojects/%v/repositories/%v/artifacts/%v",
		server,
		apiVersion,
		project,
		artifact[0],
		artifact[1],
	)

	var artifactResponse ArtifactResponse

	res, body := client(
		ClientPrt{
			Url:         url,
			Method:      "GET",
			ContentType: "application/json",
			User:        user,
			Password:    password,
			Bearer:      bearer,
		},
	)
	if res.StatusCode < 399 && res.StatusCode > 100 {

		err := json.Unmarshal([]byte(body), &artifactResponse)
		if err != nil {
			return "", errors.New("There was an error getting the artifact")
		}
	} else {

		log.Fatal("Error getting artifacts")
		return "", errors.New("There was an error getting the artifact")
	}

	return artifactResponse.Digest, nil

}

func checkArtifactSHA(clientId, tenant, server, user, password, apiVersion, project, image, sha string) bool {

	artifact, e := getArtifactSHA(clientId, tenant, server, user, password, apiVersion, project, image)

	if e != nil {
		return false
	}
	return strings.EqualFold(artifact, sha)
}

func getOidcBearer(clientId, tenant, server, user, password string) (string, error) {
	url := fmt.Sprintf("https://login.microsoftonline.com/%v/oauth2/v2.0/token",
		tenant,
	)
	var oidcTokenRequest OidcTokenRequest = OidcTokenRequest{}
	oidcTokenRequest.Client_id = clientId
	oidcTokenRequest.Password = password
	oidcTokenRequest.Username = user
	oidcTokenRequest.Response_type = "id_token"
	oidcTokenRequest.Grant_type = "password"
	oidcTokenRequest.Scope = "openid"

	res, body := client(
		ClientPrt{
			Url:         url,
			Method:      "POST",
			ContentType: "application/x-www-form-urlencoded",
			Body:        oidcTokenRequest,
		},
	)
	var oidcTokenResponse OidcTokenResponse
	if res.StatusCode < 399 && res.StatusCode > 100 {
		//log.Println(res)
		//log.Println(body)
		err := json.Unmarshal([]byte(body), &oidcTokenResponse)
		if err != nil {
			log.Fatal(err)
			return "", errors.New("There was an error getting the oidc token")
		}
	} else {

		log.Fatal("Error calling oidc endpoint")
		return "", errors.New("There was an error getting oidc token")
	}
	//log.Println(oidcTokenResponse)

	return oidcTokenResponse.IdToken, nil
}
