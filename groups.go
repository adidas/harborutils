package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

func getGroupFromName(harborServer, user, password, apiVersion, groupName string) (GroupResponse, error) {

	groups, err := listGroups(harborServer, user, password, apiVersion)

	if err == nil {
		for _, group := range groups {
			if group.GroupName == groupName {
				return group, nil
			}
		}
	}

	return GroupResponse{}, errors.New("No group found")

}

func listGroups(harborServer, user, pass, apiVersion string) (GroupListResponse, error) {
	url := fmt.Sprintf("%v/api/%vusergroups",
		harborServer,
		apiVersion,
	)
	response := GroupListResponse{}
	res, body := client(
		ClientPrt{
			Url:         url,
			Method:      "GET",
			ContentType: "application/json",
			User:        user,
			Password:    pass,
		},
	)
	if res.StatusCode < 399 && res.StatusCode > 100 {
		json.Unmarshal([]byte(body), &response)
	} else {
		log.Fatal("Error getting groups from server")

	}
	if len(response) > 0 {
		return response, nil
	}

	return GroupListResponse{}, errors.New("Groups not found")

}

func deleteGroup(harborServer, user, pass, apiVersion string, group GroupResponse) {
	url := fmt.Sprintf("%v/api/%vusergroups/%v",
		harborServer,
		apiVersion,
		group.GroupID,
	)

	res, _ := client(
		ClientPrt{
			Url:         url,
			Method:      "DELETE",
			ContentType: "application/json",
			User:        user,
			Password:    pass,
		},
	)
	if res.StatusCode > 399 || res.StatusCode < 100 {
		log.Fatal("Error getting groups from server")
	}
}
func deleteGroups(harborServer, user, pass, apiVersion string, groups GroupListResponse) {
	for _, group := range groups {
		deleteGroup(harborServer, user, pass, apiVersion, group)
	}
}

func getGroupsFromPrefix(harborServer, user, password, apiVersion, groupPrefix string) (GroupListResponse, error) {

	groups, err := listGroups(harborServer, user, password, apiVersion)

	groupsResponse := GroupListResponse{}

	if err == nil {
		for _, group := range groups {
			if strings.HasPrefix(group.GroupName, groupPrefix) {
				groupsResponse = append(groupsResponse, group)
			}
		}
		return groupsResponse, nil
	} else {
		return GroupListResponse{}, err
	}

}

func importGroups(harborServer, user, pass, apiVersion string, groups GroupListResponse){
	url := fmt.Sprintf("%v/api/%v/usergroups",
	harborServer,
	apiVersion,
	)
	
	group := GroupResponse{}

	for _, g := range groups{
		if strings.HasPrefix(group.GroupName, "CN=") {
			fmt.Println("Skip invalid group name")
		} else {
			group.GroupName = g.GroupName
			group.LdapGroupDN = g.LdapGroupDN

			res, _ := client(
				ClientPrt{
					Url:         url,
					Method:      "POST",
					ContentType: "application/json",
					User:        user,
					Password:    pass,
					Body:        group,
				},
			)
			if res.StatusCode < 100 || (res.StatusCode > 399 && res.StatusCode != 409) {
				// print status code because some groups might not be valid anymore and import should continue
				fmt.Printf("Status Code: %d while trying to import group: %v", res.StatusCode, group.GroupName)				
			} else {
				fmt.Println("Group: " + group.GroupName + " imported")
			}
		}

	}
}
