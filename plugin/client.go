package plugin

import "github.com/torlenor/abylebotter/model"

// AbyleBotterPlugin should be embedded in the Plugin to get access to the bot API
type AbyleBotterPlugin struct {
	// API exposes the plugin API of the bot.
	API     API
	Storage StorageAPI
}

// SetAPI gives the API interface to the plugin.
func (p *AbyleBotterPlugin) SetAPI(api API, storageAPI StorageAPI) {
	p.API = api
	p.Storage = storageAPI
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
