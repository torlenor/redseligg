package plugin

import "github.com/torlenor/abylebotter/model"

// Hooks are all the function the plugin has to implement to work with the bot.
// If the plugin is not interested in it, just implement it empty (which is also the default implementation).
type Hooks interface {
	// OnPost is called when a new post is made
	OnPost(model.Post)
}
