package slack

type Broadcast struct {
	OmitUsers interface{} `json:"omit_users"`
	UserID    string      `json:"user_id"`
	ChannelID string      `json:"channel_id"`
	TeamID    string      `json:"team_id"`
}

type EventStatusChange struct {
	Event string `json:"event"`
	Data  struct {
		Status string `json:"status"`
		UserID string `json:"user_id"`
	} `json:"data"`
	Broadcast Broadcast `json:"broadcast"`
	Seq       int       `json:"seq"`
}

type EventHello struct {
	Event string `json:"event"`
	Data  struct {
		ServerVersion string `json:"server_version"`
	} `json:"data"`
	Broadcast Broadcast `json:"broadcast"`
	Seq       int       `json:"seq"`
}

type EventTyping struct {
	Event string `json:"event"`
	Data  struct {
		ParentID string `json:"parent_id"`
		UserID   string `json:"user_id"`
	} `json:"data"`
	Broadcast Broadcast `json:"broadcast"`
	Seq       int       `json:"seq"`
}

type EventMessage struct {
	ClientMsgID          string `json:"client_msg_id"`
	SuppressNotification bool   `json:"suppress_notification"`
	Type                 string `json:"type"`
	Text                 string `json:"text"`
	User                 string `json:"user"`
	Team                 string `json:"team"`
	UserTeam             string `json:"user_team"`
	SourceTeam           string `json:"source_team"`
	Channel              string `json:"channel"`
	EventTs              string `json:"event_ts"`
	Ts                   string `json:"ts"`
}
