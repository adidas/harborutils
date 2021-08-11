package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func getReplicationTask(server, user, password, apiVersion string, executionID int) []ReplicationTask {
	url := fmt.Sprintf("%v/api/%vreplication/executions/%d/tasks",
		server,
		apiVersion,
		executionID,
	)

	var r []ReplicationTask

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
		json.Unmarshal([]byte(body), &r)
	} else {
		log.Fatalf("Error getting replication Execution: %d ,errorCode %d\n", executionID, res.StatusCode)
		os.Exit(1)
	}
	return r
}

func getReplicationTaksByPolicyName(server, user, password, api, policyName string) {
	rps := getReplicationPolicyByName(server, user, password, api, policyName)
	if len(rps) != 1 {
		if len(rps) == 0 {
			log.Fatalf("No policy found, policyName %s\n", policyName)
		} else {
			str := ""
			for _, rp := range rps {
				str = fmt.Sprintf("%s, %s", str, rp.Name)
			}
			log.Fatalf("Found multiple policies with the same name %s\n", str)
		}
		os.Exit(1)
	}
	policyId := rps[0].ID
	rpes := getReplicationExecutionsByPolicyId(server, user, password, api, policyId)
	for _, rpe := range rpes {
		// todo use goroutines/channels
		t := getReplicationTask(server, user, password, api, rpe.ID)
		// when it has not been able to run a replication, for example due to redis failure
		if len(t) == 0 && rpe.Status == "Failed" {
			b, _ := json.Marshal(rpe)
			fmt.Println(string(b))
		} else {
			for _, aux := range t {
				b, _ := json.Marshal(aux)
				fmt.Println(string(b))
			}
		}
	}
}
