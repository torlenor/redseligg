package platform

import (
	"context"

	"github.com/torlenor/abylebotter/plugin"
)

// Bot type interface which every Bot has to implement
type Bot interface {
	// Run the Bot (blocking)
	Run(ctx context.Context) error

	AddPlugin(plugin BotPlugin)

	GetInfo() BotInfo
}

// BotPlugin is needed to connect a Plugin to a Bot
type BotPlugin interface {
	plugin.Hooks

	SetBotPluginID(botID string, pluginID string)

	SetAPI(api plugin.API)
}

// PluginInfo contains info about one plugin
type PluginInfo struct {
	Plugin string `json:"plugin"`
	Active bool   `json:"active"`
}

// BotInfo contains info about one bot
type BotInfo struct {
	BotID    string       `json:"botId"`
	Platform string       `json:"platform"`
	Healthy  bool         `json:"healthy"`
	Plugins  []PluginInfo `json:"plugins"`
}
