package plugin

import (
	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/storage"
)

// AbyleBotterPlugin should be embedded in the Plugin to get access to the bot API
type AbyleBotterPlugin struct {
	// API exposes the plugin API of the bot.
	API     API
	Storage storage.Storage

	BotID    string
	PluginID string
}

// SetAPI gives the API interface to the plugin.
func (p *AbyleBotterPlugin) SetAPI(api API, storage storage.Storage) {
	p.API = api
	p.Storage = storage
}

// SetBotPluginID sets the plugin ID to the given value
func (p *AbyleBotterPlugin) SetBotPluginID(botID string, pluginID string) {
	p.BotID = botID
	p.PluginID = pluginID
}

// Default hook implementations (see hooks.go)

// OnRun in its default implementation
func (p *AbyleBotterPlugin) OnRun() {}

// OnStop in its default implementation
func (p *AbyleBotterPlugin) OnStop() {}

// OnPost in its default implementation.
func (p *AbyleBotterPlugin) OnPost(post model.Post) {}

// OnReactionAdded in its default implementation.
func (p *AbyleBotterPlugin) OnReactionAdded(model.Reaction) {}

// OnReactionRemoved in its default implementation.
func (p *AbyleBotterPlugin) OnReactionRemoved(model.Reaction) {}
