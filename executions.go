package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func getReplicationExecutionsByPolicyId(server, user, password, apiVersion string, policyId int) []ReplicationExecution {
	page := 0
	pageSize := 50
	baseUrl := fmt.Sprintf("%v/api/%vreplication/executions",
		server,
		apiVersion,
	)

	var replications, aux []ReplicationExecution
	var res *http.Response
	var body, url string

	for {
		page++
		url = fmt.Sprintf("%v?page=%v&page_size=%v&policy_id=%v",
			baseUrl,
			page,
			pageSize,
			policyId,
		)
		res, body = client(
			ClientPrt{
				Url:         url,
				Method:      "GET",
				ContentType: "application/json",
				User:        user,
				Password:    password,
			},
		)
		if res.StatusCode < 399 && res.StatusCode > 100 {
			json.Unmarshal([]byte(body), &aux)
			if len(aux) == 0 {
				break
			}
			replications = append(replications, aux...)

		} else {
			log.Fatal("Error getting replications")
		}

	}
	return replications
}
