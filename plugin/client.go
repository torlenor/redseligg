package plugin

import (
	"fmt"

	"github.com/torlenor/abylebotter/model"
)

// AbyleBotterPlugin should be embedded in the Plugin to get access to the bot API
type AbyleBotterPlugin struct {
	// API exposes the plugin API of the bot.
	API API

	NeededFeatures []string

	Type string

	BotID    string
	PluginID string
}

// SetAPI gives the API interface to the plugin.
func (p *AbyleBotterPlugin) SetAPI(api API) error {
	for _, f := range p.NeededFeatures {
		if !api.HasFeature(f) {
			return fmt.Errorf("Bot does not provided needed feature %s", f)
		}
	}

	p.API = api

	return nil
}

// SetBotPluginID sets the plugin ID to the given value
func (p *AbyleBotterPlugin) SetBotPluginID(botID string, pluginID string) {
	p.BotID = botID
	p.PluginID = pluginID
}

// Default hook implementations (see hooks.go)

// PluginType returns the plugin type
func (p *AbyleBotterPlugin) PluginType() string { return p.Type }

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
