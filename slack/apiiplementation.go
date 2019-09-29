package slack

import (
	"encoding/json"

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

type ConversationsListResponse struct {
	Ok               bool      `json:"ok"`
	Channels         []Channel `json:"channels"`
	ResponseMetadata struct {
		NextCursor string `json:"next_cursor"`
	} `json:"response_metadata"`
}

func (b *Bot) getConversationsList() (ConversationsListResponse, error) {
	rawResponse, err := b.apiCall("/api/conversations.list", "GET", "")
	if err != nil {
		return ConversationsListResponse{}, errors.Wrap(err, "apiCall failed")
	}

	response := ConversationsListResponse{}

	err = json.Unmarshal(rawResponse.body, &response)
	return response, err
}
