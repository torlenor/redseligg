package botinterface

import (
	"events"
	"plugins"
)

// Bot type interface
type Bot interface {
	Start(doneChannel chan struct{})
	Stop()
	Status() BotStatus

	AddPlugin(plugin plugins.Plugin)

	GetReceiveMessageChannel() chan events.ReceiveMessage
	GetSendMessageChannel() chan events.SendMessage
	GetCommandChannel() chan events.Command
}

// BotStatus gives information about the current status of the bot
type BotStatus struct {
	Running bool
	Fail    bool
	Fatal   bool
}
