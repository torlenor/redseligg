package discord

import (
	"encoding/json"
	"log"
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

func (b Bot) sendWhisper(snowflakeID string, content string) error {
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

	response, err = b.apiCall("/channels/"+channelResponseData["id"].(string)+"/messages", "POST", `{"content": "`+content+`"}`)
	if err != nil {
		return errors.Wrap(err, "apiCall failed")
	}
	if checkRateLimit(response) > 0 {
		return errors.New("sending failed (sending whisper)")
	}
	log.Printf("DiscordBot: Sent: WHISPER to UserID = %s, Content = %s", snowflakeID, content)
	return nil
}

func (b Bot) sendMessage(channelID string, content string) error {
	response, err := b.apiCall("/channels/"+channelID+"/messages", "POST", `{"content": "`+content+`"}`)
	if err != nil {
		return errors.Wrap(err, "apiCall failed")
	}
	if checkRateLimit(response) > 0 {
		return errors.New("sending failed (sending message)")
	}
	log.Printf("DiscordBot: Sent: MESSAGE to ChannelID = %s, Content = %s", channelID, content)
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
