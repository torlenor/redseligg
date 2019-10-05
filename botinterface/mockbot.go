package botinterface

import "github.com/torlenor/abylebotter/events"

type MockBot struct {
	ReceiveMessageChan chan events.ReceiveMessage
	SendMessageChan    chan events.SendMessage

	LastSendMessage events.SendMessage
}

func (b *MockBot) StartSendChannelReceiver() {
	b.LastSendMessage = <-b.SendMessageChan
}

func (b *MockBot) Reset() {
	b.LastSendMessage = events.SendMessage{Type: events.UNKNOWN, Ident: "", Content: ""}
}

func (b *MockBot) SendMessage(msg events.ReceiveMessage) {
	select {
	case b.ReceiveMessageChan <- msg:
	default:
	}
}
