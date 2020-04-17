package discord

import (
	"encoding/json"
	"fmt"
	"time"
)

type event struct {
	Op      int             `json:"op"`
	Seq     int64           `json:"s"`
	Type    string          `json:"t"`
	RawData json.RawMessage `json:"d"`
}

type hello struct {
	HeartbeatInterval int64 `json:"heartbeat_interval"`
}

type invalidSession struct {
	Op int  `json:"op"`
	D  bool `json:"d"`
}

type typingStart struct {
	UserID    string `json:"user_id"`
	Timestamp int    `json:"timestamp"`
	Member    struct {
		User struct {
			Username      string `json:"username"`
			ID            string `json:"id"`
			Discriminator string `json:"discriminator"`
			Avatar        string `json:"avatar"`
		} `json:"user"`
		Roles    []string  `json:"roles"`
		Mute     bool      `json:"mute"`
		JoinedAt time.Time `json:"joined_at"`
		Deaf     bool      `json:"deaf"`
	} `json:"member"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
}

func (t typingStart) toString() string {
	data, err := json.Marshal(t)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type presenceUpdate struct {
	User struct {
		ID string `json:"id"`
	} `json:"user"`
	Status  string      `json:"status"`
	Roles   []string    `json:"roles"`
	Nick    interface{} `json:"nick"`
	GuildID string      `json:"guild_id"`
	Game    struct {
		Type int    `json:"type"`
		Name string `json:"name"`
	} `json:"game"`
}

func (p presenceUpdate) toString() string {
	data, err := json.Marshal(p)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type messageCreate struct {
	Type            int           `json:"type"`
	Tts             bool          `json:"tts"`
	Timestamp       time.Time     `json:"timestamp"`
	Pinned          bool          `json:"pinned"`
	Nonce           interface{}   `json:"nonce"`
	Mentions        []interface{} `json:"mentions"`
	MentionRoles    []interface{} `json:"mention_roles"`
	MentionEveryone bool          `json:"mention_everyone"`
	Member          struct {
		Roles    []string  `json:"roles"`
		Mute     bool      `json:"mute"`
		JoinedAt time.Time `json:"joined_at"`
		Deaf     bool      `json:"deaf"`
	} `json:"member"`
	ID              string        `json:"id"`
	Embeds          []interface{} `json:"embeds"`
	EditedTimestamp interface{}   `json:"edited_timestamp"`
	Content         string        `json:"content"`
	ChannelID       string        `json:"channel_id"`
	Author          struct {
		Username      string      `json:"username"`
		ID            string      `json:"id"`
		Discriminator string      `json:"discriminator"`
		Bot           bool        `json:"bot"`
		Avatar        interface{} `json:"avatar"`
	} `json:"author"`
	Attachments []interface{} `json:"attachments"`
	GuildID     string        `json:"guild_id"`
}

func (mc messageCreate) toString() string {
	data, err := json.Marshal(mc)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type channel struct {
	Type                 int           `json:"type"`
	Topic                string        `json:"topic,omitempty"`
	Position             int           `json:"position"`
	PermissionOverwrites []interface{} `json:"permission_overwrites"`
	Name                 string        `json:"name"`
	LastPinTimestamp     time.Time     `json:"last_pin_timestamp,omitempty"`
	LastMessageID        string        `json:"last_message_id,omitempty"`
	ID                   string        `json:"id"`
	UserLimit            int           `json:"user_limit,omitempty"`
	ParentID             string        `json:"parent_id,omitempty"`
	Bitrate              int           `json:"bitrate,omitempty"`
}

type guildCreate struct {
	VoiceStates       []interface{} `json:"voice_states"`
	VerificationLevel int           `json:"verification_level"`
	Unavailable       bool          `json:"unavailable"`
	SystemChannelID   interface{}   `json:"system_channel_id"`
	Splash            interface{}   `json:"splash"`
	Roles             []struct {
		Position    int    `json:"position"`
		Permissions int    `json:"permissions"`
		Name        string `json:"name"`
		Mentionable bool   `json:"mentionable"`
		Managed     bool   `json:"managed"`
		ID          string `json:"id"`
		Hoist       bool   `json:"hoist"`
		Color       int    `json:"color"`
	} `json:"roles"`
	Region    string `json:"region"`
	Presences []struct {
		User struct {
			ID string `json:"id"`
		} `json:"user"`
		Status string      `json:"status"`
		Game   interface{} `json:"game"`
	} `json:"presences"`
	OwnerID  string `json:"owner_id"`
	Name     string `json:"name"`
	MfaLevel int    `json:"mfa_level"`
	Members  []struct {
		User struct {
			Username      string `json:"username"`
			ID            string `json:"id"`
			Discriminator string `json:"discriminator"`
			Avatar        string `json:"avatar"`
		} `json:"user"`
		Roles    []string  `json:"roles"`
		Mute     bool      `json:"mute"`
		JoinedAt time.Time `json:"joined_at"`
		Deaf     bool      `json:"deaf"`
		Nick     string    `json:"nick,omitempty"`
	} `json:"members"`
	MemberCount                 int           `json:"member_count"`
	Lazy                        bool          `json:"lazy"`
	Large                       bool          `json:"large"`
	JoinedAt                    time.Time     `json:"joined_at"`
	ID                          string        `json:"id"`
	Icon                        string        `json:"icon"`
	Features                    []interface{} `json:"features"`
	ExplicitContentFilter       int           `json:"explicit_content_filter"`
	Emojis                      []interface{} `json:"emojis"`
	DefaultMessageNotifications int           `json:"default_message_notifications"`
	Channels                    []channel     `json:"channels"`
	ApplicationID               interface{}   `json:"application_id"`
	AfkTimeout                  int           `json:"afk_timeout"`
	AfkChannelID                interface{}   `json:"afk_channel_id"`
}

func (gc guildCreate) toString() string {
	data, err := json.Marshal(gc)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type ready struct {
	V            int `json:"v"`
	UserSettings struct {
	} `json:"user_settings"`
	User struct {
		Verified      bool        `json:"verified"`
		Username      string      `json:"username"`
		MfaEnabled    bool        `json:"mfa_enabled"`
		ID            string      `json:"id"`
		Email         interface{} `json:"email"`
		Discriminator string      `json:"discriminator"`
		Bot           bool        `json:"bot"`
		Avatar        interface{} `json:"avatar"`
	} `json:"user"`
	SessionID       string        `json:"session_id"`
	Relationships   []interface{} `json:"relationships"`
	PrivateChannels []interface{} `json:"private_channels"`
	Presences       []interface{} `json:"presences"`
	Guilds          []struct {
		Unavailable bool   `json:"unavailable"`
		ID          string `json:"id"`
	} `json:"guilds"`
	Trace []string `json:"_trace"`
}

func (r ready) toString() string {
	data, err := json.Marshal(r)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type channelCreate struct {
	Type       int `json:"type"`
	Recipients []struct {
		Username      string `json:"username"`
		ID            string `json:"id"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar"`
	} `json:"recipients"`
	LastMessageID string `json:"last_message_id"`
	ID            string `json:"id"`
}

func (cc channelCreate) toString() string {
	data, err := json.Marshal(cc)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type messageDelete struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
}

func (md messageDelete) toString() string {
	data, err := json.Marshal(md)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type messageReactionAdd struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	Member    struct {
		User struct {
			Username      string `json:"username"`
			ID            string `json:"id"`
			Discriminator string `json:"discriminator"`
			Avatar        string `json:"avatar"`
		} `json:"user"`
		Roles       []interface{} `json:"roles"`
		Mute        bool          `json:"mute"`
		JoinedAt    time.Time     `json:"joined_at"`
		HoistedRole interface{}   `json:"hoisted_role"`
		Deaf        bool          `json:"deaf"`
	} `json:"member"`
	Emoji struct {
		Name string      `json:"name"`
		ID   interface{} `json:"id"`
	} `json:"emoji"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
}

func (mra messageReactionAdd) toString() string {
	data, err := json.Marshal(mra)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type messageReactionRemove struct {
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	Emoji     struct {
		Name     string      `json:"name"`
		ID       interface{} `json:"id"`
		Animated bool        `json:"animated"`
	} `json:"emoji"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
}

func (mrr messageReactionRemove) toString() string {
	data, err := json.Marshal(mrr)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type messageUpdate struct {
	Type            int           `json:"type"`
	Tts             bool          `json:"tts"`
	Timestamp       time.Time     `json:"timestamp"`
	Pinned          bool          `json:"pinned"`
	Nonce           interface{}   `json:"nonce"`
	Mentions        []interface{} `json:"mentions"`
	MentionRoles    []interface{} `json:"mention_roles"`
	MentionEveryone bool          `json:"mention_everyone"`
	ID              string        `json:"id"`
	Embeds          []interface{} `json:"embeds"`
	EditedTimestamp time.Time     `json:"edited_timestamp"`
	Content         string        `json:"content"`
	ChannelID       string        `json:"channel_id"`
	Author          struct {
		Username      string `json:"username"`
		ID            string `json:"id"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar"`
	} `json:"author"`
	Attachments []interface{} `json:"attachments"`
	GuildID     string        `json:"guild_id"`
}

func (mu messageUpdate) toString() string {
	data, err := json.Marshal(mu)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type channelPinsUpdate struct {
	LastPinTimestamp time.Time `json:"last_pin_timestamp"`
	ChannelID        string    `json:"channel_id"`
	GuildID          string    `json:"guild_id"`
}

func (cpu channelPinsUpdate) toString() string {
	data, err := json.Marshal(cpu)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}
	return fmt.Sprintf("%s", data)
}

type guildMemberUpdate struct {
	User struct {
		Username      string      `json:"username"`
		ID            string      `json:"id"`
		Discriminator string      `json:"discriminator"`
		Bot           bool        `json:"bot"`
		Avatar        interface{} `json:"avatar"`
	} `json:"user"`
	Roles   []string    `json:"roles"`
	Nick    interface{} `json:"nick"`
	GuildID string      `json:"guild_id"`
}

func (gmu guildMemberUpdate) toString() string {
	data, err := json.Marshal(gmu)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}

	return fmt.Sprintf("%s", data)
}

type presencesReplace struct {
	D []interface{} `json:"d"`
}

func (pr presencesReplace) toString() string {
	data, err := json.Marshal(pr)
	if err != nil {
		log.Errorln("UNHANDELED ERROR:", err)
	}

	return fmt.Sprintf("%s", data)
}
