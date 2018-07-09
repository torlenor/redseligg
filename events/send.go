package events

type SendMessage struct {
	Type    MessageType
	Ident   string
	Content string
}
