package plugin

import "github.com/torlenor/abylebotter/model"

// AbyleBotterPlugin should be embedded in the Plugin to get access to the bot API
type AbyleBotterPlugin struct {
	// API exposes the plugin API of the bot.
	API API
}

// SetAPI gives the API interface to the plugin.
func (p *AbyleBotterPlugin) SetAPI(api API) {
	p.API = api
}

// Default hook implementations (see hooks.go)

// OnPost in its default implementation
func (p *AbyleBotterPlugin) OnPost(post model.Post) {}
