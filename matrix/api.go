package matrix

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

var doHTTPRequest = func(request *http.Request) (*http.Response, error) {
	client := &http.Client{}
	return client.Do(request)
}

type api interface {
	call(path string, method string, body string, auth bool) (r []byte, e error)

	updateAuthToken(token string)

	connectToMatrixServer() error
	login(username string, password string) error
}

type matrixAPI struct {
	server    string
	authToken string
}

func (api *matrixAPI) call(path string, method string, body string, auth bool) (r []byte, e error) {

	req, err := http.NewRequest(method, api.server+"/_matrix"+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	if auth == true {
		req.Header.Add("Authorization", "Bearer "+api.authToken)
	}
	req.Header.Add("Content-Type", "application/json")

	response, err := doHTTPRequest(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(response.Body)
}

func (api *matrixAPI) connectToMatrixServer() error {
	response, err := api.call("/client/r0/join/!cJQhJDXTxLzZeuoHzw:matrix.abyle.org?access_token="+api.authToken, "POST", `{}`, false)
	if err != nil {
		return errors.Wrap(err, "apiCall failed")
	}
	log.Debugln("connectToMatrixServer() response:", string(response))
	return nil
}

func (api *matrixAPI) updateAuthToken(token string) {
	api.authToken = token
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	HomeServer  string `json:"home_server"`
	UserID      string `json:"user_id"`
	DeviceID    string `json:"device_id"`
}

func (api *matrixAPI) login(username string, password string) error {
	// get login server
	response, err := api.call("/client/r0/login", "POST", `{"type":"m.login.password", "user":"`+username+`", "password":"`+password+`"}`, false)
	if err != nil {
		return errors.Wrap(err, "apiCall failed")
	}

	var channelResponseData loginResponse
	if err := json.Unmarshal(response, &channelResponseData); err != nil {
		return errors.Wrap(err, "json unmarshal failed")
	}

	if len(channelResponseData.AccessToken) > 0 {
		api.authToken = channelResponseData.AccessToken
		return nil
	}

	return errors.New("could not login")
}
