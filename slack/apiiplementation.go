package slack

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type Channel struct {
	ID                      string        `json:"id"`
	Name                    string        `json:"name"`
	IsChannel               bool          `json:"is_channel"`
	IsGroup                 bool          `json:"is_group"`
	IsIm                    bool          `json:"is_im"`
	Created                 int           `json:"created"`
	IsArchived              bool          `json:"is_archived"`
	IsGeneral               bool          `json:"is_general"`
	Unlinked                int           `json:"unlinked"`
	NameNormalized          string        `json:"name_normalized"`
	IsShared                bool          `json:"is_shared"`
	ParentConversation      interface{}   `json:"parent_conversation"`
	Creator                 string        `json:"creator"`
	IsExtShared             bool          `json:"is_ext_shared"`
	IsOrgShared             bool          `json:"is_org_shared"`
	SharedTeamIds           []string      `json:"shared_team_ids"`
	PendingShared           []interface{} `json:"pending_shared"`
	PendingConnectedTeamIds []interface{} `json:"pending_connected_team_ids"`
	IsPendingExtShared      bool          `json:"is_pending_ext_shared"`
	IsMember                bool          `json:"is_member"`
	IsPrivate               bool          `json:"is_private"`
	IsMpim                  bool          `json:"is_mpim"`
	Topic                   struct {
		Value   string `json:"value"`
		Creator string `json:"creator"`
		LastSet int    `json:"last_set"`
	} `json:"topic"`
	Purpose struct {
		Value   string `json:"value"`
		Creator string `json:"creator"`
		LastSet int    `json:"last_set"`
	} `json:"purpose"`
	PreviousNames []interface{} `json:"previous_names"`
	NumMembers    int           `json:"num_members"`
}

type Channels []Channel

type ConversationsListResponse struct {
	Ok               bool     `json:"ok"`
	Channels         Channels `json:"channels"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}

func (b *Bot) getConversationsList(cursor string) (ConversationsListResponse, error) {
	args := ""
	if len(cursor) > 0 {
		args = "cursor=" + cursor
	}
	rawResponse, err := b.apiCall("/api/conversations.list", "GET", args, "")
	if err != nil {
		return ConversationsListResponse{}, errors.Wrap(err, "apiCall failed")
	}

	response := ConversationsListResponse{}
	err = json.Unmarshal(rawResponse.body, &response)

	if err == nil && !response.Ok {
		return ConversationsListResponse{}, fmt.Errorf("Error getting conversations.list: Received not OK")
	} else if err != nil {
		return ConversationsListResponse{}, err
	}

	return response, nil
}

func (b *Bot) getConversations() (Channels, error) {
	conversationsList, err := b.getConversationsList("")
	if err != nil {
		return Channels{}, fmt.Errorf("Error getting conversations.list: %s", err)
	}

	channels := conversationsList.Channels

	nextCursor := conversationsList.ResponseMetadata.NextCursor
	for len(nextCursor) != 0 {
		conversationsList, err := b.getConversationsList(nextCursor)
		if err != nil {
			return Channels{}, fmt.Errorf("Error getting conversations.list: %s", err)
		}
		channels = append(channels, conversationsList.Channels...)
		nextCursor = conversationsList.ResponseMetadata.NextCursor
	}

	return channels, nil
}

type User struct {
	ID       string `json:"id"`
	TeamID   string `json:"team_id"`
	Name     string `json:"name"`
	Deleted  bool   `json:"deleted"`
	Color    string `json:"color"`
	RealName string `json:"real_name"`
	Tz       string `json:"tz"`
	TzLabel  string `json:"tz_label"`
	TzOffset int    `json:"tz_offset"`
	Profile  struct {
		AvatarHash            string `json:"avatar_hash"`
		StatusText            string `json:"status_text"`
		StatusEmoji           string `json:"status_emoji"`
		RealName              string `json:"real_name"`
		DisplayName           string `json:"display_name"`
		RealNameNormalized    string `json:"real_name_normalized"`
		DisplayNameNormalized string `json:"display_name_normalized"`
		Email                 string `json:"email"`
		Image24               string `json:"image_24"`
		Image32               string `json:"image_32"`
		Image48               string `json:"image_48"`
		Image72               string `json:"image_72"`
		Image192              string `json:"image_192"`
		Image512              string `json:"image_512"`
		Team                  string `json:"team"`
		FirstName             string `json:"first_name"`
		LastName              string `json:"last_name"`
		Title                 string `json:"title"`
		Phone                 string `json:"phone"`
		Skype                 string `json:"skype"`
	} `json:"profile"`
	IsAdmin           bool `json:"is_admin"`
	IsOwner           bool `json:"is_owner"`
	IsPrimaryOwner    bool `json:"is_primary_owner"`
	IsRestricted      bool `json:"is_restricted"`
	IsUltraRestricted bool `json:"is_ultra_restricted"`
	IsBot             bool `json:"is_bot"`
	Updated           int  `json:"updated"`
	IsAppUser         bool `json:"is_app_user"`
	Has2Fa            bool `json:"has_2fa"`
}

type Users []User

type UserListResponse struct {
	Ok               bool  `json:"ok"`
	Members          Users `json:"members"`
	CacheTs          int   `json:"cache_ts"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}

func (b *Bot) getUserList(cursor string) (UserListResponse, error) {
	args := ""
	if len(cursor) > 0 {
		args = "cursor=" + cursor
	}
	rawResponse, err := b.apiCall("/api/users.list", "GET", args, "")
	if err != nil {
		return UserListResponse{}, errors.Wrap(err, "apiCall failed")
	}

	response := UserListResponse{}
	err = json.Unmarshal(rawResponse.body, &response)

	if err == nil && !response.Ok {
		return UserListResponse{}, fmt.Errorf("Error getting users.list: Received not OK")
	} else if err != nil {
		return UserListResponse{}, err
	}

	return response, nil
}

func (b *Bot) getUsers() (Users, error) {
	usersList, err := b.getUserList("")
	if err != nil {
		return Users{}, fmt.Errorf("Error getting users.list: %s", err)
	}

	users := usersList.Members

	nextCursor := usersList.ResponseMetadata.NextCursor
	for len(nextCursor) != 0 {
		conversationsList, err := b.getConversationsList(nextCursor)
		if err != nil {
			return Users{}, fmt.Errorf("Error getting users.list: %s", err)
		}
		users = append(users, usersList.Members...)
		nextCursor = conversationsList.ResponseMetadata.NextCursor
	}

	return users, nil
}
