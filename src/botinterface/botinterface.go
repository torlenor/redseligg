package botinterface

import "../events"

// Bot type interface
type Bot interface {
	Start(doneChannel chan struct{})
	Stop()

	GetReceiveMessageChannel() chan events.ReceiveMessage
	GetSendMessageChannel() chan events.SendMessage
	GetCommandChannel() chan events.Command
}
