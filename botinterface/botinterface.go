package botinterface

// Bot type interface
type Bot interface {
	Start()
	Stop()
	Status() BotStatus
}

// BotStatus gives information about the current status of the bot
type BotStatus struct {
	Running bool
	Fail    bool
	Fatal   bool
}
