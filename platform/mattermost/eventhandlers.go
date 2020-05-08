package mattermost

import (
	"encoding/json"

	"github.com/torlenor/redseligg/model"
)

type post struct {
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
	var posted eventPosted

	if err := json.Unmarshal(data, &posted); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	var post post

	if err := json.Unmarshal([]byte(posted.Data.Post), &post); err != nil {
		b.log.Errorln("UNHANDLED ERROR: ", err)
		return
	}

	b.log.Printf("%s", data)

	isPrivate := false
	if posted.Data.ChannelType == "D" {
		isPrivate = true
	}

	var userName string
	user, err := b.getUserByID(post.UserID)
	if err != nil {
		userName = user.Username
	}

	receiveMessage := model.Post{ServerID: b.config.Server, User: model.User{Name: userName, ID: post.UserID}, ChannelID: post.ChannelID, Content: post.Message, IsPrivate: isPrivate}
	for _, plugin := range b.plugins {
		plugin.OnPost(receiveMessage)
	}

}
