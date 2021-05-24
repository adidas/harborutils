package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func getReplicationPolicy(server, user, password, apiVersion string, policyId int) ReplicationPolicy {
	url := fmt.Sprintf("%v/api/%vreplication/policies/%d",
		server,
		apiVersion,
		policyId,
	)

	var rp ReplicationPolicy

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
		json.Unmarshal([]byte(body), &rp)
	} else {
		log.Fatal("Error getting replication policy ", policyId)
		os.Exit(1)

	}
	return rp
}

func updateReplication(server, user, password, apiVersion, image, tag string, policyId int, rp ReplicationPolicy) {
	// image := strings.Split(resource, ":")[0]
	// tag := strings.Split(resource, ":")[1]
	rp.Filters = []RPFilter{{"name", image}, RPFilter{"tag", tag}, RPFilter{"resource", "image"}}
	url := fmt.Sprintf("%v/api/%vreplication/policies/%d",
		server,
		apiVersion,
		policyId,
	)
	res, _ := client(
		ClientPrt{
			Url:         url,
			Method:      "PUT",
			ContentType: "application/json",
			User:        user,
			Password:    password,
			Body:        rp,
		},
	)
	if res.StatusCode < 100 || (res.StatusCode > 399 && res.StatusCode != 409) {
		fmt.Printf("Status Code: %d while trying to update Replication Policy: %d, image: %s, tag: %s\n", res.StatusCode, policyId, image, tag)
	} else {
		fmt.Printf("ReplicationPolicy updated: %d, image: %s, tag: %s\n", policyId, image, tag)
	}
}

func startReplication(server, user, password, apiVersion string, policyId int) {
	url := fmt.Sprintf("%v/api/%vreplication/executions",
		server,
		apiVersion,
	)
	r := StartReplicationExecution{policyId}
	res, _ := client(
		ClientPrt{
			Url:         url,
			Method:      "POST",
			ContentType: "application/json",
			User:        user,
			Password:    password,
			Body:        r,
		},
	)
	if res.StatusCode < 100 || (res.StatusCode > 399 && res.StatusCode != 409) {
		fmt.Printf("Fail Start Execution\n")
	} else {
		fmt.Printf("Start Execution\n")
	}
}

func getLastExecution(server, user, password, apiVersion string, policyId int) ReplicationExecution {
	url := fmt.Sprintf("%v/api/%vreplication/executions?policy_id=%d&page=1&page_size=15",
		server,
		apiVersion,
		policyId,
	)

	var r []ReplicationExecution

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
		log.Fatal("Error getting replication Executin ", policyId)
		os.Exit(1)

	}
	return r[0]
}

func getReplicationExecution(server, user, password, apiVersion string, executionID int) ReplicationExecution {
	url := fmt.Sprintf("%v/api/%vreplication/executions/%d",
		server,
		apiVersion,
		executionID,
	)

	var r ReplicationExecution

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
		log.Fatal("Error getting replication Executin ", executionID)
		os.Exit(1)
	}
	return r
}

func waitFinishReplication(server, user, password, api, resource string, policyId int) {
	execution := getLastExecution(server, user, password, api, policyId)
	for execution.Status == "InProgress" {
		time.Sleep(10 * time.Second)
		execution = getReplicationExecution(server, user, password, api, execution.ID)
	}
	if execution.Status != "Succeed" {
		fmt.Printf("ERROR fail replication, %s execution: %+v\n", resource, execution)
	} else {
		fmt.Printf("Replication finished: %s\n", resource)
	}
}

func compactReplication(server, user, password, api string, logs []AuditLog) map[string]string {
	// var c map[string]string
	c := make(map[string]string)

	for _, log := range logs {
		image := strings.Split(log.Resource, ":")[0]
		tag := strings.Split(log.Resource, ":")[1]
		if v, ok := c[image]; ok {
			c[image] = v + tag + ","
		} else {
			c[image] = tag + ","
		}
	}
	return c
}

func replication(server, user, password, api string, startAt, finishAt time.Time, policyId int) {
	rp := getReplicationPolicy(server, user, password, api, policyId)
	fmt.Println(rp)
	logs := listAuditLogs(server, user, password, api, startAt, finishAt)
	// for _, log := range logs {
	// 	fmt.Println(log)
	// 	updateReplication(server, user, password, api, log.Resource, policyId, rp)
	// 	startReplication(server, user, password, api, policyId)
	// 	waitFinishReplication(server, user, password, api, log.Resource, policyId)
	// 	// maybe we can remove the repeated ones
	// }
	c := compactReplication(server, user, password, api, logs)
	l := len(c)
	i := 0
	for image, tag := range c {
		fmt.Printf("Start replication %d/%d\n", i, l)
		updateReplication(server, user, password, api, image, "{"+tag+"}", policyId, rp)
		startReplication(server, user, password, api, policyId)
		waitFinishReplication(server, user, password, api, image+":"+tag, policyId)
		// maybe we can remove the repeated ones
		i++
	}

}
