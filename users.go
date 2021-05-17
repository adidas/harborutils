package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
)

func getUserFromUsername(harborServer, user, password, apiVersion, userName string) (UserResponse, error) {
	url := fmt.Sprintf("%v/api/%vusers?username=%v",
		harborServer,
		apiVersion,
		userName,
	)
	return getUser(url, user, password)
}

func getUserFromEmail(harborServer, user, password, apiVersion, email string) (UserResponse, error) {
	url := fmt.Sprintf("%v/api/%vusers?email=%v",
		harborServer,
		apiVersion,
		email,
	)

	return getUser(url, user, password)
}

func getUser(url, user, pass string) (UserResponse, error) {
	response := []UserResponse{}
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
		log.Fatal("Error getting Ldap information")
	}
	if len(response) > 0 {
		return response[0], nil
	}

	return UserResponse{}, errors.New("User not found")

}

func getSourceUsers(harborServer, user, pass, apiVersion string) (UserListResponse, error) {
	var users UserListResponse
	cnt := 1
	for {
		url := fmt.Sprintf("%v/api/%vusers?page=%v",
			harborServer,
			apiVersion,
			cnt)

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
			var u UserListResponse
			json.Unmarshal([]byte(body), &u)
			users = append(users, u...)
		} else {
			log.Fatal("Error getting users")
		}
		
		l := res.Header.Get("Link")
		if strings.Contains(l, "next") {
			cnt++			
		} else {
			break
		}
	}	
	if len(users) > 0 {
		return users, nil
	}
	return users, errors.New("Error getting users")

}

func importLdapUser(harborServer, user, pass, apiVersion string, userImport UserImport) {
	url := fmt.Sprintf("%v/api/%vldap/users/import",
		harborServer,
		apiVersion,
	)
	
	res, _ := client(
		ClientPrt{
			Url:         url,
			Method:      "POST",
			ContentType: "application/json",
			User:        user,
			Password:    pass,
			Body:        userImport,
		},
	)
	
	if res.StatusCode < 100 || (res.StatusCode > 399 && res.StatusCode != 409) {
		//print status code as if user not found to import it returns 404
		log.Println(res.StatusCode)
		log.Fatal("Error importing ldap user")
	}
	
}
