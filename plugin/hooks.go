package plugin

import "github.com/torlenor/abylebotter/model"

// Hooks are all the function the plugin has to implement to work with the bot.
// If the plugin is not interested in it, just implement it empty (which is also the default implementation).
type Hooks interface {
	// PluginType returns the plugin type
	PluginType() string
	// OnRun is called when the bot is operational and ready to run
	OnRun()
	// OnStop is called when the bot is shutting down
	OnStop()
	// OnPost is called when a new post is made
	OnPost(model.Post)
	// OnReactionAdded is called when a reaction to posted message is received. This can be, e.g., an emoji.
	OnReactionAdded(model.Reaction)
	// OnReactionRemoved is called when a reaction is removed from a posted message. This can be, e.g., an emoji.
	OnReactionRemoved(model.Reaction)
}
