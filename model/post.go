package model

// Post is a event of either incoming or outgoing messages
type Post struct {
	ServerID string // ServerID is a unique ID on which the Bot identifies a server
	Server   string // [optional] ServerID is the clear text name of a Server

	ChannelID string // ChannelID is a unique ID on which the Bot identifies the channel
	Channel   string // [optional] Channel is a clear text name of a Channel in which the message was seen

	User User // User is the sender/receiver

	Content string

	IsPrivate bool // IsPrivate indicates it is a whisper or similar (depending on the Bot)
}

// MessageIdentifier is a unique identifier for a message on a platform
type MessageIdentifier struct {
	ID      string
	Channel string
}

// PostResponse contains infos about posted messages
type PostResponse struct {
	PostedMessageIdent MessageIdentifier
}
