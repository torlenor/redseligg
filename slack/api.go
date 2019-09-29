package slack

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func (b *Bot) apiRunner(path string, method string, args string, body string) (*apiResponse, error) {
	finished := false
	tries := 0
	for !finished {
		tries++
		if tries > 3 {
			return nil, errors.New("API call still failing after 3 tries, giving up")
		}

		response, err := b.apiCall(path, method, args, body)
		if err != nil {
			return response, errors.Wrap(err, "apiCall failed")
		}

		b.log.Printf("SlackBot: API Call %s %s %s finished", path, method, body)
		return response, nil
	}

	return nil, nil
}

// RtmConnectResponse contains WebSocket Message Server URL and limited information about the team
type RtmConnectResponse struct {
	Ok   bool   `json:"ok"`
	URL  string `json:"url"`
	Team struct {
		ID     string `json:"id"`
		Name   string `json:"name"`
		Domain string `json:"domain"`
	} `json:"team"`
	Self struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"self"`
}

// RtmConnect returns a WebSocket Message Server URL and limited information about the team
func (b *Bot) RtmConnect() (RtmConnectResponse, error) {
	// get login server
	response, err := b.apiCall("/api/rtm.connect", "GET", "", "")
	if err != nil {
		return RtmConnectResponse{}, errors.Wrap(err, "apiCall failed")
	}

	rtmConnectResponse := RtmConnectResponse{}

	err = json.Unmarshal(response.body, &rtmConnectResponse)
	return rtmConnectResponse, err
}

type apiResponse struct {
	header     http.Header
	body       []byte
	statusCode int
}

func (b *Bot) apiCall(path string, method string, args string, body string) (*apiResponse, error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, "https://slack.com"+path+"?token="+b.config.Token+args, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

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
