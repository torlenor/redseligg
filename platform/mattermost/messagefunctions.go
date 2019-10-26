package mattermost

import (
	"encoding/json"

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
	// It is a known channel
	if _, ok := b.knownChannelIDs[userID]; ok {
		return b.sendMessage(userID, content)
	}

	// It is not a known Channel so maybe it is a userID
	response, err := b.apiRunner("/api/v4/channels/direct", "POST", `[
			"`+b.MeUser.ID+`",
			"`+userID+`"
			]`)
	if err != nil && response.statusCode > 200 {
		return err
	}
	var channel channelData
	err = json.Unmarshal(response.body, &channel)
	if err != nil {
		return err
	}
	b.addKnownChannel(channel)
	return b.sendMessage(channel.ID, content)
}
