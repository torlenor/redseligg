package mattermost

type broadcast struct {
	OmitUsers interface{} `json:"omit_users"`
	UserID    string      `json:"user_id"`
	ChannelID string      `json:"channel_id"`
	TeamID    string      `json:"team_id"`
}

type eventStatusChange struct {
	Event string `json:"event"`
	Data  struct {
		Status string `json:"status"`
		UserID string `json:"user_id"`
	} `json:"data"`
	Broadcast broadcast `json:"broadcast"`
	Seq       int       `json:"seq"`
}

type eventHello struct {
	Event string `json:"event"`
	Data  struct {
		ServerVersion string `json:"server_version"`
	} `json:"data"`
	Broadcast broadcast `json:"broadcast"`
	Seq       int       `json:"seq"`
}

type eventTyping struct {
	Event string `json:"event"`
	Data  struct {
		ParentID string `json:"parent_id"`
		UserID   string `json:"user_id"`
	} `json:"data"`
	Broadcast broadcast `json:"broadcast"`
	Seq       int       `json:"seq"`
}

type eventPosted struct {
	Event string `json:"event"`
	Data  struct {
		ChannelDisplayName string `json:"channel_display_name"`
		ChannelName        string `json:"channel_name"`
		ChannelType        string `json:"channel_type"`
		Post               string `json:"post"`
		SenderName         string `json:"sender_name"`
		TeamID             string `json:"team_id"`
	} `json:"data"`
	Broadcast struct {
		OmitUsers interface{} `json:"omit_users"`
		UserID    string      `json:"user_id"`
		ChannelID string      `json:"channel_id"`
		TeamID    string      `json:"team_id"`
	} `json:"broadcast"`
	Seq int `json:"seq"`
}
