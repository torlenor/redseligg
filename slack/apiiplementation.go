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
