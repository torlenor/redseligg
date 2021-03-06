package slack

type eventMessage struct {
	ClientMsgID          string `json:"client_msg_id"`
	SuppressNotification bool   `json:"suppress_notification"`
	Type                 string `json:"type"`
	Subtype              string `json:"subtype"`
	Text                 string `json:"text"`
	User                 string `json:"user"`
	Team                 string `json:"team"`
	UserTeam             string `json:"user_team"`
	SourceTeam           string `json:"source_team"`
	Channel              string `json:"channel"`
	EventTs              string `json:"event_ts"`
	Ts                   string `json:"ts"`
}

type eventUserTyping struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	User    string `json:"user"`
}

type eventDesktopNotification struct {
	Type            string `json:"type"`
	Title           string `json:"title"`
	Subtitle        string `json:"subtitle"`
	Msg             string `json:"msg"`
	Ts              string `json:"ts"`
	Content         string `json:"content"`
	Channel         string `json:"channel"`
	LaunchURI       string `json:"launchUri"`
	AvatarImage     string `json:"avatarImage"`
	SsbFilename     string `json:"ssbFilename"`
	ImageURI        string `json:"imageUri"`
	IsShared        bool   `json:"is_shared"`
	IsChannelInvite bool   `json:"is_channel_invite"`
	EventTs         string `json:"event_ts"`
}

type eventChannelCreated struct {
	Type    string `json:"type"`
	Channel struct {
		ID             string `json:"id"`
		IsChannel      bool   `json:"is_channel"`
		Name           string `json:"name"`
		NameNormalized string `json:"name_normalized"`
		Created        int    `json:"created"`
		Creator        string `json:"creator"`
		IsShared       bool   `json:"is_shared"`
		IsOrgShared    bool   `json:"is_org_shared"`
	} `json:"channel"`
	EventTs string `json:"event_ts"`
}

type eventChannel struct {
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
	IsFrozen                bool          `json:"is_frozen"`
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
	LastRead                string        `json:"last_read"`
	Latest                  struct {
		User    string `json:"user"`
		Type    string `json:"type"`
		Subtype string `json:"subtype"`
		Ts      string `json:"ts"`
		Text    string `json:"text"`
		Inviter string `json:"inviter"`
	} `json:"latest"`
	UnreadCount        int      `json:"unread_count"`
	UnreadCountDisplay int      `json:"unread_count_display"`
	IsOpen             bool     `json:"is_open"`
	Members            []string `json:"members"`
	Topic              struct {
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
	Priority      int           `json:"priority"`
}

type eventChannelJoined struct {
	Type    string       `json:"type"`
	Channel eventChannel `json:"channel"`
}

type eventChannelLeft struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	ActorID string `json:"actor_id"`
	EventTs string `json:"event_ts"`
}

type eventChannelDeleted struct {
	Type    string `json:"type"`
	Channel string `json:"channel"`
	ActorID string `json:"actor_id"`
	EventTs string `json:"event_ts"`
}

type eventMemberJoinedChannel struct {
	Type        string `json:"type"`
	User        string `json:"user"`
	Channel     string `json:"channel"`
	ChannelType string `json:"channel_type"`
	Team        string `json:"team"`
	Inviter     string `json:"inviter"`
	EventTs     string `json:"event_ts"`
	Ts          string `json:"ts"`
}

type eventGroupJoined struct {
	Type    string       `json:"type"`
	Channel eventChannel `json:"channel"`
}

type eventAck struct {
	Ok      bool   `json:"ok"`
	ReplyTo int    `json:"reply_to"`
	Ts      string `json:"ts"`
	Text    string `json:"text"`
}

type eventUser struct {
	Type    string `json:"type"`
	User    user   `json:"user"`
	CacheTs int    `json:"cache_ts"`
	EventTs string `json:"event_ts"`
}

type eventDnDUpdatedUser struct {
	Type      string `json:"type"`
	User      string `json:"user"`
	DndStatus struct {
		DndEnabled     bool `json:"dnd_enabled"`
		NextDndStartTs int  `json:"next_dnd_start_ts"`
		NextDndEndTs   int  `json:"next_dnd_end_ts"`
	} `json:"dnd_status"`
	EventTs string `json:"event_ts"`
}

type eventIMCreated struct {
	Type    string `json:"type"`
	User    string `json:"user"`
	Channel struct {
		ID                 string      `json:"id"`
		Created            int         `json:"created"`
		IsFrozen           bool        `json:"is_frozen"`
		IsArchived         bool        `json:"is_archived"`
		IsIm               bool        `json:"is_im"`
		IsOrgShared        bool        `json:"is_org_shared"`
		User               string      `json:"user"`
		LastRead           string      `json:"last_read"`
		Latest             interface{} `json:"latest"`
		UnreadCount        int         `json:"unread_count"`
		UnreadCountDisplay int         `json:"unread_count_display"`
		IsOpen             bool        `json:"is_open"`
	} `json:"channel"`
	EventTs string `json:"event_ts"`
}

type eventReactionAddedOrRemoved struct {
	Type string `json:"type"`
	User string `json:"user"`
	Item struct {
		Type    string `json:"type"`
		Channel string `json:"channel"`
		Ts      string `json:"ts"`
	} `json:"item"`
	Reaction string `json:"reaction"`
	ItemUser string `json:"item_user"`
	EventTs  string `json:"event_ts"`
	Ts       string `json:"ts"`
}
