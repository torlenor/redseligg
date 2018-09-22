package botinterface

import (
	"plugins"
)

// Bot type interface
type Bot interface {
	Start(doneChannel chan struct{})
	Stop()
	Status() BotStatus

	AddPlugin(plugin plugins.Plugin)
}

// BotStatus gives information about the current status of the bot
type BotStatus struct {
	Running bool
	Fail    bool
	Fatal   bool
}
