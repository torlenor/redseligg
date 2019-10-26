package platform

import "github.com/torlenor/abylebotter/plugin"

// Bot type interface which every Bot has to implement
type Bot interface {
	Start()
	Stop()

	AddPlugin(plugin BotPlugin)
}

// BotPlugin is needed to connect a Plugin to a Bot
type BotPlugin interface {
	plugin.Hooks
	SetAPI(api plugin.API)
}
