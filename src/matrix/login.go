package matrix

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type loginResponse struct {
	AccessToken string `json:"access_token"`
	HomeServer  string `json:"home_server"`
	UserID      string `json:"user_id"`
	DeviceID    string `json:"device_id"`
}

func (b *Bot) connectToMatrixServer() error {
	response, err := b.apiCall("/client/r0/join/!cJQhJDXTxLzZeuoHzw:matrix.abyle.org?access_token="+b.token, "POST", `{}`, false)
	if err != nil {
		return errors.Wrap(err, "apiCall failed")
	}
	log.Debugln("connectToMatrixServer() response:", string(response))
	return nil
}

func (b *Bot) login(username string, password string) (string, error) {
	// get login server
	response, err := b.apiCall("/client/r0/login", "POST", `{"type":"m.login.password", "user":"`+username+`", "password":"`+password+`"}`, false)
	if err != nil {
		return "", errors.Wrap(err, "apiCall failed")
	}

	log.Debugln("login() response:", string(response))

	var channelResponseData loginResponse
	if err := json.Unmarshal(response, &channelResponseData); err != nil {
		return "", errors.Wrap(err, "json unmarshal failed")
	}

	if len(channelResponseData.AccessToken) > 0 {
		return channelResponseData.AccessToken, nil
	}

	return string(""), errors.New("could not login")
}
