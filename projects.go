package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func listProjects(server, user, password, api string) ProjectListResponse {
	page := 0
	pageSize := 10
	baseUrl := fmt.Sprintf("%v/api/%vprojects",
		server,
		api,
	)

	var projects, aux ProjectListResponse
	var res *http.Response
	var body, url string

	for {
		page++
		url = fmt.Sprintf("%v?page=%v&page_size=%v",
			baseUrl,
			page,
			pageSize,
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
			projects = append(projects, aux...)

		} else {
			log.Fatal("Error getting projects")
		}

	}

	return projects
}

func getProject(server, user, password, name, apiVersion string) (ProjectResponse, error) {
	url := fmt.Sprintf("%v/api/%vprojects?name=%v",
		server,
		apiVersion,
		name,
	)

	var projects ProjectListResponse

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
		json.Unmarshal([]byte(body), &projects)
	} else {
		log.Fatal("Error getting projects")
	}
	if len(projects) > 0 {
		return projects[0], nil
	}
	return ProjectResponse{}, errors.New("No project Found")
}

func listMembers(server, user, password, apiVersion string, projectID int) (MemberListResponse, MemberListResponse) {
	url := fmt.Sprintf("%v/api/%vprojects/%v/members",
		server,
		apiVersion,
		projectID,
	)
	groups := MemberListResponse{}
	users := MemberListResponse{}

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
		response := MemberListResponse{}
		json.Unmarshal([]byte(body), &response)

		for _, member := range response {
			if member.EntityType == "g" {
				groups = append(groups, member)
			} else {
				users = append(users, member)
			}
		}
	} else {
		log.Fatal("Error getting Members")
	}
	return groups, users
}

func addMember(server, user, password, apiVersion string, projectID int, member AddMember) {
	url := fmt.Sprintf("%v/api/%vprojects/%v/members",
		server,
		apiVersion,
		projectID,
	)

	res, _ := client(
		ClientPrt{
			Url:         url,
			Method:      "POST",
			ContentType: "application/json",
			Body:        member,
			User:        user,
			Password:    password,
		},
	)
	if res.StatusCode < 100 || (res.StatusCode > 399 && res.StatusCode != 409) {
		log.Fatal("Error adding Members")
	}

}

func syncProjectGrants(project ProjectResponse) {
	groups, users := listMembers(harborServer, harborUser, harborPassword, harborAPIVersion, project.ProjectId)
	var newUser UserResponse
	var err error
	newProject, err := getProject(harborServerTarget, harborUserTarget, harborPasswordTarget, project.Name, harborAPIVersionTarget)
	if err == nil {
		for _, group := range groups {
			newGroup, err := getGroupFromName(harborServerTarget, harborUserTarget, harborPasswordTarget, harborAPIVersionTarget, group.EntityName)
			if err == nil {
				//fmt.Println("Group ", newGroup.GroupName, " with role ", group.RoleId, " Found")

				newMember := AddMember{}
				newMember.RoleID = group.RoleId
				newMember.MemberGroup.ID = newGroup.GroupID
				newMember.MemberGroup.GroupName = newGroup.GroupName
				addMember(harborServerTarget, harborUserTarget, harborPasswordTarget, harborAPIVersionTarget, newProject.ProjectId, newMember)

			} else {
				log.Println("Group not found in target server: ", group.EntityName)
			}
		}
		for _, user := range users {
			u, _ := getUserFromUsername(harborServer, harborUser, harborPassword, harborAPIVersion, user.EntityName)

			if strings.Contains(u.Username, "@") || strings.Contains(u.Username, "hlp_") {
				newUser, err = getUserFromUsername(harborServerTarget, harborUserTarget, harborPasswordTarget, harborAPIVersionTarget, u.Username)

			} else {
				if strings.Contains(u.Email, "@placeholder.com") {
					newUser, err = getUserFromEmail(harborServerTarget, harborUserTarget, harborPasswordTarget, harborAPIVersionTarget, u.Username+"@emea.adsint.biz")
				} else {
					newUser, err = getUserFromEmail(harborServerTarget, harborUserTarget, harborPasswordTarget, harborAPIVersionTarget, u.Email)
				}
			}
			if err == nil {
				newMember := AddMember{}
				newMember.RoleID = user.RoleId
				newMember.MemberUser.UserID = newUser.UserID
				newMember.MemberUser.Username = newUser.Username
				addMember(harborServerTarget, harborUserTarget, harborPasswordTarget, harborAPIVersionTarget, newProject.ProjectId, newMember)

			} else {
				log.Println("User not found in target server: ", u.Username)
			}

		}
	} else {
		log.Println("Project not found in target server: ", project.Name)
	}

}

func syncProjectLabels(project ProjectResponse) {
	labels, _ := getProjectLabels(harborServer, harborUser, harborPassword, harborAPIVersion, project.ProjectId)
	var err error
	newProject, err := getProject(harborServerTarget, harborUserTarget, harborPasswordTarget, project.Name, harborAPIVersionTarget)
	if err == nil {
		for _, label := range labels {
			newLabel := LabelResponse{}
			newLabel.Description = label.Description
			newLabel.Name = label.Name
			newLabel.Scope = label.Scope
			newLabel.ProjectID = newProject.ProjectId
			addLabel(harborServerTarget, harborUserTarget, harborPasswordTarget, harborAPIVersionTarget, newProject.ProjectId, newLabel)
		}

	} else {
		log.Println("Project not found in target server: ", project.Name)
	}

}
