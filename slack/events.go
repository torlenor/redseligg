package slack

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
