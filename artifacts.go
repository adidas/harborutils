package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

func getArtifactSHA(server, user, password, apiVersion, project, image string) (string, error) {
	artifact := strings.Split(image, ":")
	if len(artifact) < 2 {
		artifact = append(artifact, "latest")
	}

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

func checkArtifactSHA(server, user, password, apiVersion, project, image, sha string) bool {

	artifact, e := getArtifactSHA(server, user, password, apiVersion, project, image)

	if e != nil {
		return false
	}
	return strings.EqualFold(artifact, sha)
}
