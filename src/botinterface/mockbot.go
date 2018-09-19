package botinterface

import "events"

type MockBot struct {
	ReceiveMessageChan chan events.ReceiveMessage
	SendMessageChan    chan events.SendMessage
	CommandChan        chan events.Command

	LastSendMessage     events.SendMessage
	LastReceivedCommand events.Command
}

func (b *MockBot) StartSendChannelReceiver() {
	b.LastSendMessage = <-b.SendMessageChan
}

func (b *MockBot) StartCommandChannelReceiver() {
	b.LastReceivedCommand = <-b.CommandChan
}

func (b *MockBot) Reset() {
	b.LastSendMessage = events.SendMessage{Type: events.UNKNOWN, Ident: "", Content: ""}
	b.LastReceivedCommand = events.Command{Command: "", Payload: ""}
}

func (b *MockBot) SendMessage(msg events.ReceiveMessage) {
	select {
	case b.ReceiveMessageChan <- msg:
	default:
	}
}
