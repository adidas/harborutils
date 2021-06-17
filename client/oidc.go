package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func GetOidcBearer(clientId, tenant, server, user, password string) (string, error) {
	url := fmt.Sprintf("https://login.microsoftonline.com/%v/oauth2/v2.0/token",
		tenant,
	)
	var oidcTokenRequest OidcTokenRequest = OidcTokenRequest{}
	oidcTokenRequest.Client_id = clientId
	oidcTokenRequest.Password = password
	oidcTokenRequest.Username = user
	oidcTokenRequest.Response_type = "id_token"
	oidcTokenRequest.Grant_type = "password"
	oidcTokenRequest.Scope = "openid"

	res, body := client(
		ClientPrt{
			Url:         url,
			Method:      "POST",
			ContentType: "application/x-www-form-urlencoded",
			Body:        oidcTokenRequest,
		},
	)
	var oidcTokenResponse OidcTokenResponse
	if res.StatusCode < 399 && res.StatusCode > 100 {
		//log.Println(res)
		//log.Println(body)
		err := json.Unmarshal([]byte(body), &oidcTokenResponse)
		if err != nil {
			log.Fatal(err)
			return "", errors.New("There was an error getting the oidc token")
		}
	} else {
		return "", errors.New("There was an error getting oidc token")
	}
	//log.Println(oidcTokenResponse)

	return oidcTokenResponse.IdToken, nil
}
