package events

// SendMessage is used to notify about messages to send
type SendMessage struct {
	Type    MessageType
	Ident   string
	Content string
}
