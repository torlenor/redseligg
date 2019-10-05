package events

// SendMessage is used to notify about messages to send
type SendMessage struct {
	Type MessageType

	ChannelID string // ChannelID is a unique ID on which the Bot identifies the channel on which this message was seen
	UserID    string // UserID is a unique ID on which the Bot identifies the User which sent the message

	Content string
}
