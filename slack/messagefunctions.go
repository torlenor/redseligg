package slack

import "github.com/pkg/errors"

type SendMessage struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func (b *Bot) sendMessage(channelID string, content string) error {

	msg := SendMessage{
		ID:      1, // TODO needs to increase
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
	// // It is a known channel
	// if _, ok := b.knownChannelIDs[userID]; ok {
	return b.sendMessage(userID, content)
	// } else {
	// 	// It is not a known Channel so maybe it is a userID
	// 	response, err := b.apiRunner("/api/v4/channels/direct", "POST", `[
	// 		"`+userID+`"
	// 		]`)
	// 	if err != nil && response.statusCode > 200 {
	// 		return err
	// 	}
	// 	var channel Channel
	// 	err = json.Unmarshal(response.body, &channel)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	b.addKnownChannel(channel)
	// 	return b.sendMessage(channel.ID, content)
	// }
}
