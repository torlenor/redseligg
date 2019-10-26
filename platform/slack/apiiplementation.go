package slack

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
)

type channel struct {
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

type channels []channel

type conversationsListResponse struct {
	Ok               bool     `json:"ok"`
	Channels         channels `json:"channels"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}

func (b *Bot) getConversationsList(cursor string) (conversationsListResponse, error) {
	args := ""
	if len(cursor) > 0 {
		args = "cursor=" + cursor
	}
	rawResponse, err := b.apiCall("/api/conversations.list", "GET", args, "")
	if err != nil {
		return conversationsListResponse{}, errors.Wrap(err, "apiCall failed")
	}

	response := conversationsListResponse{}
	err = json.Unmarshal(rawResponse.body, &response)

	if err == nil && !response.Ok {
		return conversationsListResponse{}, fmt.Errorf("Error getting conversations.list: Received not OK")
	} else if err != nil {
		return conversationsListResponse{}, err
	}

	return response, nil
}

func (b *Bot) getConversations() (channels, error) {
	conversationsList, err := b.getConversationsList("")
	if err != nil {
		return channels{}, fmt.Errorf("Error getting conversations.list: %s", err)
	}

	c := conversationsList.Channels

	nextCursor := conversationsList.ResponseMetadata.NextCursor
	for len(nextCursor) != 0 {
		conversationsList, err := b.getConversationsList(nextCursor)
		if err != nil {
			return channels{}, fmt.Errorf("Error getting conversations.list: %s", err)
		}
		c = append(c, conversationsList.Channels...)
		nextCursor = conversationsList.ResponseMetadata.NextCursor
	}

	return c, nil
}

type userListResponse struct {
	Ok               bool  `json:"ok"`
	Members          users `json:"members"`
	CacheTs          int   `json:"cache_ts"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}

func (b *Bot) getUserList(cursor string) (userListResponse, error) {
	args := ""
	if len(cursor) > 0 {
		args = "cursor=" + cursor
	}
	rawResponse, err := b.apiCall("/api/users.list", "GET", args, "")
	if err != nil {
		return userListResponse{}, errors.Wrap(err, "apiCall failed")
	}

	response := userListResponse{}
	err = json.Unmarshal(rawResponse.body, &response)

	if err == nil && !response.Ok {
		return userListResponse{}, fmt.Errorf("Error getting users.list: Received not OK")
	} else if err != nil {
		return userListResponse{}, err
	}

	return response, nil
}

func (b *Bot) getUsers() (users, error) {
	usersList, err := b.getUserList("")
	if err != nil {
		return users{}, fmt.Errorf("Error getting users.list: %s", err)
	}

	us := usersList.Members

	nextCursor := usersList.ResponseMetadata.NextCursor
	for len(nextCursor) != 0 {
		conversationsList, err := b.getConversationsList(nextCursor)
		if err != nil {
			return users{}, fmt.Errorf("Error getting users.list: %s", err)
		}
		us = append(us, usersList.Members...)
		nextCursor = conversationsList.ResponseMetadata.NextCursor
	}

	return us, nil
}
