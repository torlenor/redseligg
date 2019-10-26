package discord

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type rateLimitResponse struct {
	Global     bool   `json:"global"`
	Message    string `json:"message"`
	RetryAfter int    `json:"retry_after"`
}

type successfulSentResponse struct {
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

func (b *Bot) sendWhisper(snowflakeID string, content string) error {
	response, err := b.apiCall("/users/@me/channels", "POST", `{"recipient_id": "`+snowflakeID+`"}`)
	if err != nil {
		return errors.Wrap(err, "apiCall failed")
	}
	if checkRateLimit(response) > 0 {
		return errors.New("sending failed (create channel)")
	}

	var channelResponseData map[string]interface{}
	if err := json.Unmarshal(response, &channelResponseData); err != nil {
		return errors.Wrap(err, "json unmarshal failed")
	}

	var channelID string
	if id, ok := channelResponseData["id"].(string); ok {
		channelID = id
	} else {
		return errors.Wrap(err, "no valid channel id found")
	}

	err = b.messageRunner(channelID, content)
	if err == nil {
		b.stats.whispersSent++
	}

	return err
}

func (b *Bot) sendMessage(receiver string, content string) error {
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
			return errors.New("Unknown Guild " + guild + ". Cannot send message")
		}

		for _, entry := range b.guilds[guildID].Channels {
			if entry.ID == channel || entry.Name == channel {
				channelID = entry.ID
				break
			}
		}
	}

	err := b.messageRunner(channelID, content)
	if err == nil {
		b.stats.messagesSent++
	}

	return err
}

func (b *Bot) messageRunner(channelID string, content string) error {
	finished := false
	tries := 0
	for !finished {
		tries++
		if tries > 3 {
			return errors.New("Message sending still failing after 3 tries, giving up")
		}

		response, err := b.apiCall("/channels/"+channelID+"/messages", "POST", `{"content": "`+content+`"}`)
		if err != nil {
			return errors.Wrap(err, "apiCall failed")
		}
		var retryAfter int
		if retryAfter = checkRateLimit(response); retryAfter > 0 {
			log.Warn("Sending failed because we are rate limited. Trying to resend after: " + strconv.Itoa(retryAfter))
			time.Sleep(time.Duration(retryAfter) * time.Millisecond)
			continue
		}
		log.Printf("DiscordBot: Sent: MESSAGE to ChannelID = %s, Content = %s", channelID, content)
		finished = true
	}

	return nil
}

func checkRateLimit(response []byte) int {
	var rateLimited rateLimitResponse
	err := json.Unmarshal(response, &rateLimited)
	if err != nil {
		return 0
	}
	return rateLimited.RetryAfter
}