package plugin

import (
	"fmt"

	"github.com/torlenor/redseligg/model"
)

// RedseliggPlugin should be embedded in the Plugin to get access to the bot API
type RedseliggPlugin struct {
	// API exposes the plugin API of the bot.
	API API

	NeededFeatures []string

	Type string

	BotID    string
	PluginID string
}

// SetAPI gives the API interface to the plugin.
func (p *RedseliggPlugin) SetAPI(api API) error {
	for _, f := range p.NeededFeatures {
		if !api.HasFeature(f) {
			return fmt.Errorf("Bot does not provided needed feature %s", f)
		}
	}

	p.API = api

	return nil
}

// SetBotPluginID sets the plugin ID to the given value
func (p *RedseliggPlugin) SetBotPluginID(botID string, pluginID string) {
	p.BotID = botID
	p.PluginID = pluginID
}

// Default hook implementations (see hooks.go)

// PluginType returns the plugin type
func (p *RedseliggPlugin) PluginType() string { return p.Type }

// OnRun in its default implementation
func (p *RedseliggPlugin) OnRun() {}

// OnStop in its default implementation
func (p *RedseliggPlugin) OnStop() {}

// OnPost in its default implementation.
func (p *RedseliggPlugin) OnPost(post model.Post) {}

// OnReactionAdded in its default implementation.
func (p *RedseliggPlugin) OnReactionAdded(model.Reaction) {}

// OnReactionRemoved in its default implementation.
func (p *RedseliggPlugin) OnReactionRemoved(model.Reaction) {}
