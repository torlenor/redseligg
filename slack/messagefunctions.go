package slack

import "github.com/pkg/errors"

type sendMessage struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func (b *Bot) sendMessage(channelID string, content string) error {

	msg := sendMessage{
		ID:      b.idProvider.Get(),
		Type:    "message",
		Channel: channelID,
		Text:    content,
	}

	err := b.ws.SendJSONMessage(msg)
	if err != nil {
		return errors.New("Sending Message failed: " + err.Error())
	}

	return nil
}

func (b *Bot) sendWhisper(userID string, content string) error {
	return b.sendMessage(userID, content)
}
