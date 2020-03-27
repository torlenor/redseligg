package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func realAPICall(serverURL string, path string, method string, body string) (r []byte, status int, e error) {
	request, err := http.NewRequest(method, serverURL+path, strings.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	request.Header.Set("Content-type", "application/json")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, 0, err
	}

	rbody, err := ioutil.ReadAll(response.Body)

	return rbody, response.StatusCode, err
}
