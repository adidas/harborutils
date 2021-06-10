package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	////Uncomment for debug
	//"io"
	//"os"
)

//Check if something was wrong
func check(e error, callback func()) {
	if e != nil {
		msg := fmt.Sprintf(
			"Something went wrong, ERROR: %s.",
			e,
		)
		log.Fatal(msg)
	} else {
		callback()
	}
}

//Bitbucket REST API call function
func client(c ClientPrt) (*http.Response, string) {
	client := http.Client{}
	buf := bytes.Buffer{}
	data := url.Values{}

	//Create form if needed otherwise create json body

	switch c.Body.(type) {
	case AddMember:
		b, _ := json.Marshal(c.Body.(AddMember))
		buf.WriteString(string(b))
	case LabelResponse:
		b, _ := json.Marshal(c.Body.(LabelResponse))
		buf.WriteString(string(b))
	case UserImport:
		b, _ := json.Marshal(c.Body.(UserImport))
		buf.WriteString(string(b))
	case GroupResponse:
		b, _ := json.Marshal(c.Body.(GroupResponse))
		buf.WriteString(string(b))
	case RobotResponse:
		b, _ := json.Marshal(c.Body.(RobotResponse))
		buf.WriteString(string(b))
	case ReplicationPolicy:
		b, _ := json.Marshal(c.Body.(ReplicationPolicy))
		buf.WriteString(string(b))
	case StartReplicationExecution:
		b, _ := json.Marshal(c.Body.(StartReplicationExecution))
		buf.WriteString(string(b))
	case OidcTokenRequest:

		data.Set("client_id", (c.Body.(OidcTokenRequest)).Client_id)
		data.Set("response_type", (c.Body.(OidcTokenRequest)).Response_type)
		data.Set("grant_type", (c.Body.(OidcTokenRequest)).Grant_type)
		data.Set("scope", (c.Body.(OidcTokenRequest)).Scope)
		data.Set("username", (c.Body.(OidcTokenRequest)).Username)
		data.Set("password", (c.Body.(OidcTokenRequest)).Password)
		buf.WriteString(data.Encode())
	}

	//fmt.Printf("+ curl -X %s %v <<< %v\n", c.Method, c.Url, buf.String())

	req, err := http.NewRequest(c.Method, c.Url, strings.NewReader(data.Encode()))
	if err != nil {

		log.Fatal(err)
	}
	req.Header.Add("Content-Type", c.ContentType)

	if c.Bearer == "" {
		switch c.Body.(type) {
		case OidcTokenRequest:
		default:
			req.SetBasicAuth(c.User, c.Password)

		}

	} else {

		req.Header.Set("authorization", c.Bearer)
	}

	//log.Println(req)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	bodyString := string(bodyBytes)

	////Uncomment for debug
	// if _, err := io.Copy(os.Stderr, resp.Body); err != nil {
	// 	log.Fatal(err)
	// }
	//log.Println(bodyString)

	return resp, bodyString

}
