package slack

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type User struct {
	ID            string `json:"id"`
	CreateAt      int    `json:"create_at"`
	UpdateAt      int    `json:"update_at"`
	DeleteAt      int    `json:"delete_at"`
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Nickname      string `json:"nickname"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	AuthService   string `json:"auth_service"`
	Roles         string `json:"roles"`
	Locale        string `json:"locale"`
	NotifyProps   struct {
		Email        string `json:"email"`
		Push         string `json:"push"`
		Desktop      string `json:"desktop"`
		DesktopSound string `json:"desktop_sound"`
		MentionKeys  string `json:"mention_keys"`
		Channel      string `json:"channel"`
		FirstName    string `json:"first_name"`
	} `json:"notify_props"`
	Props struct {
	} `json:"props"`
	LastPasswordUpdate int  `json:"last_password_update"`
	LastPictureUpdate  int  `json:"last_picture_update"`
	FailedAttempts     int  `json:"failed_attempts"`
	MfaActive          bool `json:"mfa_active"`
}

type Users []User

func (b *Bot) getUserByID(userID string) (*User, error) {
	if val, ok := b.KnownUsers[userID]; ok {
		return &val, nil
	}

	response, err := b.apiRunner("/api/v4/users/ids", "POST", `[
		"`+userID+`"
		]`)
	if err != nil && response.statusCode > 200 {
		return nil, err
	}
	var users Users
	err = json.Unmarshal(response.body, &users)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, errors.New("Could not find user with UserID: " + userID)
	}

	b.addKnownUser(users[0])

	return &users[0], nil
}

type Channel struct {
	ID            string `json:"id"`
	CreateAt      int    `json:"create_at"`
	UpdateAt      int    `json:"update_at"`
	DeleteAt      int    `json:"delete_at"`
	TeamID        string `json:"team_id"`
	Type          string `json:"type"`
	DisplayName   string `json:"display_name"`
	Name          string `json:"name"`
	Header        string `json:"header"`
	Purpose       string `json:"purpose"`
	LastPostAt    int    `json:"last_post_at"`
	TotalMsgCount int    `json:"total_msg_count"`
	ExtraUpdateAt int    `json:"extra_update_at"`
	CreatorID     string `json:"creator_id"`
}

func (b *Bot) getChannelByID(channelID string) (*Channel, error) {
	if val, ok := b.KnownChannels[channelID]; ok {
		return &val, nil
	}

	response, err := b.apiRunner("/api/v4/channels/"+channelID, "GET", "")
	if err != nil && response.statusCode > 200 {
		return nil, err
	}
	var channel Channel
	err = json.Unmarshal(response.body, &channel)
	if err != nil {
		return nil, err
	}

	b.addKnownChannel(channel)

	return &channel, nil
}
