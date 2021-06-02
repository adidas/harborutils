package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

func listAuditLogs(server, user, password, api string, startAt, finishAt time.Time) []AuditLog {
	page := 0
	pageSize := 100
	baseUrl := fmt.Sprintf("%v/api/%vaudit-logs?q=op_time=[%s~%s],resource_type=artifact,operation=create",
		server,
		api,
		startAt.Format("2006-01-02T15:04:05"),
		finishAt.Format("2006-01-02T15:04:05"),
	)

	var logs, aux []AuditLog
	var res *http.Response
	var body, url string

	for {
		page++
		url = fmt.Sprintf("%v&page=%v&page_size=%v",
			baseUrl,
			page,
			pageSize,
		)
		fmt.Println(url)
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
			logs = append(logs, aux...)

		} else {
			log.Fatal("Error getting logs")
		}
	}

	return logs
}
