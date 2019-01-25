package mattermost

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func (b *Bot) apiRunner(path string, method string, body string) (*apiResponse, error) {
	finished := false
	tries := 0
	for !finished {
		tries++
		if tries > 3 {
			return nil, errors.New("API call still failing after 3 tries, giving up")
		}

		response, err := b.apiCall(path, method, body)
		if err != nil {
			return response, errors.Wrap(err, "apiCall failed")
		}

		b.log.Printf("MattermostBot: API Call %s %s %s finished", path, method, body)
		return response, nil
	}

	return nil, nil
}

type loginResponse struct {
	AccessToken string `json:"access_token"`
	HomeServer  string `json:"home_server"`
	UserID      string `json:"user_id"`
	DeviceID    string `json:"device_id"`
}

func (b *Bot) login() error {
	// get login server
	response, err := b.apiCall("/api/v4/users/login", "POST", `{"login_id":"`+b.config.Username+`","password":"`+b.config.Password+`"}`)
	if err != nil {
		return errors.Wrap(err, "apiCall failed")
	}

	if val, ok := response.header["Token"]; ok {
		if len(val) > 0 {
			b.token = val[0]
		}
	} else {
		return errors.New("could not login: Response: " + string(response.body))
	}

	err = json.Unmarshal(response.body, &b.MeUser)
	return err
}

type apiResponse struct {
	header     http.Header
	body       []byte
	statusCode int
}

func (b *Bot) apiCall(path string, method string, body string) (*apiResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, b.config.Server+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+b.token)
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &apiResponse{
		body:       responseBody,
		header:     response.Header,
		statusCode: response.StatusCode,
	}, nil
}
