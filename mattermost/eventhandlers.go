package mattermost

import (
	"encoding/json"

	"github.com/torlenor/abylebotter/events"
)

type Post struct {
	ID         string `json:"id"`
	CreateAt   int64  `json:"create_at"`
	UpdateAt   int64  `json:"update_at"`
	EditAt     int    `json:"edit_at"`
	DeleteAt   int    `json:"delete_at"`
	IsPinned   bool   `json:"is_pinned"`
	UserID     string `json:"user_id"`
	ChannelID  string `json:"channel_id"`
	RootID     string `json:"root_id"`
	ParentID   string `json:"parent_id"`
	OriginalID string `json:"original_id"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	Props      struct {
	} `json:"props"`
	Hashtags      string `json:"hashtags"`
	PendingPostID string `json:"pending_post_id"`
}

func (b *Bot) handleEventPosted(data []byte) {
	var posted EventPosted

	if err := json.Unmarshal(data, &posted); err != nil {
		b.log.Errorln("UNHANDELED ERROR: ", err)
		return
	}

	var post Post

	if err := json.Unmarshal([]byte(posted.Data.Post), &post); err != nil {
		b.log.Errorln("UNHANDELED ERROR: ", err)
		return
	}

	b.log.Printf("%s", data)

	var messageType events.MessageType
	var ident string

	if posted.Data.ChannelType == "D" {
		messageType = events.WHISPER
		ident = post.UserID
	} else {
		messageType = events.MESSAGE
		ident = post.ChannelID
	}

	receiveMessage := events.ReceiveMessage{Type: messageType, ChannelID: ident, Content: post.Message}

	_, _ = b.getUserByID(post.UserID)

	for plugin, pluginChannel := range b.receivers {
		b.log.Debugln("Notifying plugin", plugin.GetName(), "about new message/whisper")
		select {
		case pluginChannel <- receiveMessage:
		default:
		}
	}

}
