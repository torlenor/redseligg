package mattermost

import (
	"github.com/pkg/errors"
)

func (b *Bot) sendMessage(channelID string, content string) error {

	_, err := b.apiRunner("/api/v4/posts", "POST", `{
		"channel_id": "`+channelID+`",
		"message": "`+content+`"
		}`)

	if err != nil {
		return errors.New("Sending Message failed: " + err.Error())
	}

	return nil
}

func (b *Bot) sendWhisper(userID string, content string) error {
	return nil
}

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
