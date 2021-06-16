package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func GetArtifactSHA(server, apiVersion, idtoken, user, password, project, image string) (string, error) {
	artifact := strings.Split(image, ":")
	if len(artifact) < 2 {
		artifact = append(artifact, "latest")
	}
	// artifact[0] = url.QueryEscape(artifact[0])
	artifact[0] = strings.ReplaceAll(artifact[0], "/", "%252F")

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
			Bearer:      bearer,
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

		return "", errors.New(body)
	}

	return artifactResponse.Digest, nil

}

func CheckArtifactSHA(server, apiVersion, idtoken, user, password, project, image, sha string) bool {

	// artifact, e := GetArtifactSHA(server, apiVersion, idtoken, project, image)
	artifact, e := GetArtifactSHA(server, apiVersion, idtoken, user, password, project, image)

	if e != nil {
		return false
	}
	return strings.EqualFold(artifact, sha)
}
