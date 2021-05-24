package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
	}

	// fmt.Printf("+ curl -X %s %v <<< %v\n", c.Method, c.Url, buf.String())
	req, err := http.NewRequest(c.Method, c.Url, &buf)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Content-Type", c.ContentType)
	req.SetBasicAuth(c.User, c.Password)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println(resp.Status)

	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	resp.Body = ioutil.NopCloser(bytes.NewReader(bodyBytes))
	bodyString := string(bodyBytes)

	////Uncomment for debug
	// if _, err := io.Copy(os.Stderr, resp.Body); err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println()

	return resp, bodyString

}
