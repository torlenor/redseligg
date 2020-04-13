package discord

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/torlenor/abylebotter/model"
)

type rateLimitResponse struct {
	Global     bool   `json:"global"`
	Message    string `json:"message"`
	RetryAfter int    `json:"retry_after"`
}

type messageObject struct {
	Nonce           interface{}   `json:"nonce"`
	Attachments     []interface{} `json:"attachments"`
	Tts             bool          `json:"tts"`
	Embeds          []interface{} `json:"embeds"`
	Timestamp       time.Time     `json:"timestamp"`
	MentionEveryone bool          `json:"mention_everyone"`
	ID              string        `json:"id"`
	Pinned          bool          `json:"pinned"`
	EditedTimestamp interface{}   `json:"edited_timestamp"`
	Author          struct {
		Username      string      `json:"username"`
		Discriminator string      `json:"discriminator"`
		Bot           bool        `json:"bot"`
		ID            string      `json:"id"`
		Avatar        interface{} `json:"avatar"`
	} `json:"author"`
	MentionRoles []interface{} `json:"mention_roles"`
	Content      string        `json:"content"`
	ChannelID    string        `json:"channel_id"`
	Mentions     []interface{} `json:"mentions"`
	Type         int           `json:"type"`
}

func (b *Bot) sendWhisper(snowflakeID string, content string) (messageObject, error) {
	response, _, err := b.apiCall("/users/@me/channels", "POST", `{"recipient_id": "`+snowflakeID+`"}`)
	if err != nil {
		return messageObject{}, errors.Wrap(err, "apiCall failed")
	}
	if checkRateLimit(response) > 0 {
		return messageObject{}, errors.New("sending failed (create channel)")
	}

	var channelResponseData map[string]interface{}
	if err := json.Unmarshal(response, &channelResponseData); err != nil {
		return messageObject{}, errors.Wrap(err, "json unmarshal failed")
	}

	var channelID string
	if id, ok := channelResponseData["id"].(string); ok {
		channelID = id
	} else {
		return messageObject{}, errors.Wrap(err, "no valid channel id found")
	}

	mo, err := b.messageRunner(channelID, content)
	if err == nil {
		b.stats.whispersSent++
	}

	return mo, err
}

func (b *Bot) sendMessage(receiver string, content string) (messageObject, error) {
	var channelID string

	splitString := strings.Split(receiver, "#")
	if len(splitString) != 2 {
		log.Errorf("Error decoding '%s' into Guild/Server and Channel. Format must be Guild#Channel. We will try using it as a channelID", receiver)
		channelID = receiver
	} else {
		guild := splitString[0]
		channel := strings.ToLower(splitString[1])

		var guildID string

		if _, ok := b.guilds[guild]; ok {
			guildID = guild
		} else if val, ok := b.guildNameToID[guild]; ok {
			guildID = val
		} else {
			return messageObject{}, errors.New("Unknown Guild " + guild + ". Cannot send message")
		}

		for _, entry := range b.guilds[guildID].Channels {
			if entry.ID == channel || entry.Name == channel {
				channelID = entry.ID
				break
			}
		}
	}

	mo, err := b.messageRunner(channelID, content)
	if err == nil {
		b.stats.messagesSent++
	}

	return mo, err
}

func (b *Bot) updateMessage(messageIdent model.MessageIdentifier, content string) (messageObject, error) {
	mo, err := b.updateRunner(messageIdent.Channel, messageIdent.ID, content)
	return mo, err
}

func (b *Bot) deleteMessage(messageIdent model.MessageIdentifier) error {
	return b.deleteRunner(messageIdent.Channel, messageIdent.ID)
}

func (b *Bot) messageRunner(channelID string, content string) (messageObject, error) {
	for tries := 0; tries < 4; tries++ {
		if tries > 3 {
			return messageObject{}, errors.New("Message sending still failing after 3 tries, giving up")
		}

		response, _, err := b.apiCall("/channels/"+channelID+"/messages", "POST", `{"content": "`+convertMessageFromAbyleBotter(content)+`"}`)
		if err != nil {
			return messageObject{}, errors.Wrap(err, "apiCall failed")
		}
		var retryAfter int
		if retryAfter = checkRateLimit(response); retryAfter > 0 {
			log.Warn("Sending failed because we are rate limited. Trying to resend after: " + strconv.Itoa(retryAfter))
			time.Sleep(time.Duration(retryAfter) * time.Millisecond)
			continue
		}
		log.Debugf("DiscordBot: Sent: MESSAGE to ChannelID = %s, Content = %s", channelID, content)
		return getMessageObject(response)
	}

	return messageObject{}, nil
}

func (b *Bot) updateRunner(channelID, messageID, content string) (messageObject, error) {
	for tries := 0; tries < 4; tries++ {
		if tries > 3 {
			return messageObject{}, errors.New("Message update still failing after 3 tries, giving up")
		}

		response, _, err := b.apiCall("/channels/"+channelID+"/messages/"+messageID, "PATCH", `{"content": "`+convertMessageFromAbyleBotter(content)+`"}`)
		if err != nil {
			return messageObject{}, errors.Wrap(err, "apiCall failed")
		}
		var retryAfter int
		if retryAfter = checkRateLimit(response); retryAfter > 0 {
			log.Warn("Sending failed because we are rate limited. Trying to resend after: " + strconv.Itoa(retryAfter))
			time.Sleep(time.Duration(retryAfter) * time.Millisecond)
			continue
		}
		log.Debugf("DiscordBot: Update MESSAGE in ChannelID = %s, MessageID = %s, Content = %s", channelID, messageID, content)
		return getMessageObject(response)
	}

	return messageObject{}, nil
}

func (b *Bot) deleteRunner(channelID string, messageID string) error {
	for tries := 0; tries < 4; tries++ {
		if tries > 3 {
			return errors.New("Message delete still failing after 3 tries, giving up")
		}

		response, statusCode, err := b.apiCall("/channels/"+channelID+"/messages/"+messageID, "DELETE", "")
		if err != nil {
			return errors.Wrap(err, "apiCall failed")
		}
		var retryAfter int
		if retryAfter = checkRateLimit(response); retryAfter > 0 {
			log.Warn("Deleting failed because we are rate limited. Trying to resend after: " + strconv.Itoa(retryAfter))
			time.Sleep(time.Duration(retryAfter) * time.Millisecond)
			continue
		}
		log.Debugf("DiscordBot: Deleted MESSAGE from ChannelID = %s, MessageID = %s", channelID, messageID)
		if statusCode != 204 {
			return errors.Wrap(err, "error deleting message")
		}
		return nil
	}

	return nil
}

func getMessageObject(response []byte) (messageObject, error) {
	var messageObject messageObject
	err := json.Unmarshal(response, &messageObject)
	return messageObject, err
}

func checkRateLimit(response []byte) int {
	var rateLimited rateLimitResponse
	err := json.Unmarshal(response, &rateLimited)
	if err != nil {
		return 0
	}
	return rateLimited.RetryAfter
}
