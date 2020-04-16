package webclient

import (
	"io/ioutil"
	"net/http"
	"strings"
)

// Client implements all the features needed for a real web client
type Client struct {
	path                string
	authorizationHeader string
	defaultContentType  string
}

// New creates a new web client
func New(path, authorizationHeader, defaultContentType string) *Client {
	return &Client{
		path:                path,
		authorizationHeader: authorizationHeader,
		defaultContentType:  defaultContentType,
	}
}

// Used for injection in unit tests of the webclient
var doHTTPRequest = func(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(request)
}

// Call an API endpoint with the given path, method and body
func (c *Client) Call(path string, method string, body string) (APIResponse, error) {
	req, err := http.NewRequest(method, c.path+path, strings.NewReader(body))
	if err != nil {
		return APIResponse{}, err
	}

	req.Header.Add("Authorization", c.authorizationHeader)
	req.Header.Add("Content-Type", c.defaultContentType)

	response, err := doHTTPRequest(req)
	if err != nil {
		return APIResponse{}, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return APIResponse{}, err
	}

	return APIResponse{
		Body:       responseBody,
		Header:     response.Header,
		StatusCode: response.StatusCode,
	}, nil
}
