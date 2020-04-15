package discord

import (
	"io/ioutil"
	"net/http"
	"strings"
)

func (b *Bot) apiCall(path string, method string, body string) (r []byte, statusCode int, e error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, "https://discordapp.com/api"+path, strings.NewReader(body))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Add("Authorization", "Bot "+b.token)
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}

	responseBody, err := ioutil.ReadAll(response.Body)

	return responseBody, response.StatusCode, err
}
