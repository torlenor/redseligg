package matrix

import (
	"encoding/json"

	"fmt"

	"github.com/mitchellh/mapstructure"
)

type room struct {
	RoomID              string
	UnreadNotifications struct {
	} `json:"unread_notifications"`
	Timeline struct {
		Limited   bool   `json:"limited"`
		PrevBatch string `json:"prev_batch"`
		Events    []struct {
			OriginServerTs int64  `json:"origin_server_ts"`
			Sender         string `json:"sender"`
			EventID        string `json:"event_id"`
			Unsigned       struct {
				Age int `json:"age"`
			} `json:"unsigned"`
			Content struct {
				Body    string `json:"body"`
				Msgtype string `json:"msgtype"`
			} `json:"content"`
			Type string `json:"type"`
		} `json:"events"`
	} `json:"timeline"`
	State struct {
		Events []struct {
			OriginServerTs int64  `json:"origin_server_ts"`
			Sender         string `json:"sender"`
			EventID        string `json:"event_id"`
			Unsigned       struct {
				Age int64 `json:"age"`
			} `json:"unsigned"`
			StateKey string `json:"state_key"`
			Content  struct {
				JoinRule string `json:"join_rule"`
				Name     string `json:"name"`
			} `json:"content"`
			Type       string `json:"type"`
			Membership string `json:"membership,omitempty"`
		} `json:"events"`
	} `json:"state"`
	InviteState struct {
		Events []struct {
			Content struct {
				Membership  string `json:"membership"`
				AvatarURL   string `json:"avatar_url"`
				Displayname string `json:"displayname"`
			} `json:"content"`
			Type           string `json:"type"`
			Sender         string `json:"sender"`
			StateKey       string `json:"state_key"`
			OriginServerTs int64  `json:"origin_server_ts,omitempty"`
			EventID        string `json:"event_id,omitempty"`
			Unsigned       struct {
				PrevContent struct {
					Membership string `json:"membership"`
				} `json:"prev_content"`
				PrevSender    string `json:"prev_sender"`
				ReplacesState string `json:"replaces_state"`
				Age           int    `json:"age"`
			} `json:"unsigned,omitempty"`
			Membership string `json:"membership,omitempty"`
		} `json:"events"`
	} `json:"invite_state"`
	AccountData struct {
		Events []interface{} `json:"events"`
	} `json:"account_data"`
}

type syncResponse struct {
	NextBatch string `json:"next_batch"`
	Rooms     struct {
		Invite []room
		Join   []room
		Leave  []room
	}
}

func (sr syncResponse) toString() string {
	data, err := json.Marshal(sr)
	if err != nil {
		log.Println("UNHANDELED ERROR: ", err)
	}
	return fmt.Sprintf("%s", data)
}

func parseRoom(roomID string, roomContent map[string]interface{}) (room, error) {
	var r room
	err := mapstructure.Decode(roomContent, &r)
	if err != nil {
		log.Println("Error decoding room", err)
		return r, err
	}
	r.RoomID = roomID

	return r, nil
}

func syncResponseFromMap(content map[string]interface{}) (syncResponse, error) {
	var sr syncResponse

	if val, ok := content["next_batch"].(string); ok {
		sr.NextBatch = val
	}

	if roomGroups, ok := content["rooms"].(map[string]interface{}); ok {
		if rooms, ok := roomGroups["leave"].(map[string]interface{}); ok {
			for roomID, value := range rooms {
				if roomEntry, ok := value.(map[string]interface{}); ok {
					r, err := parseRoom(roomID, roomEntry)
					if err != nil {
						log.Println("Error decoding room", err)
					}
					sr.Rooms.Leave = append(sr.Rooms.Leave, r)
				}
			}
		}

		if rooms, ok := roomGroups["join"].(map[string]interface{}); ok {
			for roomID, value := range rooms {
				if roomEntry, ok := value.(map[string]interface{}); ok {
					r, err := parseRoom(roomID, roomEntry)
					if err != nil {
						log.Println("Error decoding room", err)
					}
					sr.Rooms.Join = append(sr.Rooms.Join, r)
				}
			}
		}

		if rooms, ok := roomGroups["invite"].(map[string]interface{}); ok {
			for roomID, value := range rooms {
				if roomEntry, ok := value.(map[string]interface{}); ok {
					log.Println(roomEntry)
					r, err := parseRoom(roomID, roomEntry)
					if err != nil {
						log.Println("Error decoding room", err)
					}
					sr.Rooms.Invite = append(sr.Rooms.Invite, r)
					log.Println(r)
				}
			}
		}
	}

	return sr, nil
}
