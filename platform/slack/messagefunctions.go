package slack

import "github.com/pkg/errors"

type sendMessage struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func (b *Bot) sendMessage(channelID string, content string) (genericChatResponse, error) {
	return b.chatPostMessage(channelID, content)
	// return messagePostResponse{}, b.sendMessageViaRTM(channelID, content)
}

func (b *Bot) sendWhisper(userID string, content string) (genericChatResponse, error) {
	return b.sendMessage(userID, content)
}

func (b *Bot) sendMessageViaRTM(channelID string, content string) error {
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
