package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func getProjectLabels(server, user, password, apiVersion string, projectID int) (LabelListResponse, error) {

	url := fmt.Sprintf("%v/api/%vlabels?scope=p&project_id=%v",
		server,
		apiVersion,
		projectID,
	)

	var labels LabelListResponse

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
		json.Unmarshal([]byte(body), &labels)
	} else {
		log.Fatal("Error getting labels")
	}
	if len(labels) > 0 {
		return labels, nil
	}
	return LabelListResponse{}, errors.New("No labels Found")
}

func addLabel(server, user, password, apiVersion string, projectID int, label LabelResponse) {
	url := fmt.Sprintf("%v/api/%vlabels",
		server,
		apiVersion,
	)

	res, _ := client(
		ClientPrt{
			Url:         url,
			Method:      "POST",
			ContentType: "application/json",
			Body:        label,
			User:        user,
			Password:    password,
		},
	)
	if res.StatusCode < 100 || (res.StatusCode > 399 && res.StatusCode != 409) {
		log.Fatal("Error adding labels")
	}

}
