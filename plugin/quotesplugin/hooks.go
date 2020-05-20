package quotesplugin

import (
	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/utils"
)

// OnRun is called when the platform is ready
func (p *QuotesPlugin) OnRun() {
	p.API.RegisterCommand(p, "quote")
	p.API.RegisterCommand(p, "quoteadd")
	p.API.RegisterCommand(p, "quotehelp")
	p.API.RegisterCommand(p, "quoteremove")

	p.storage = p.getStorage()
	if p.storage == nil {
		p.API.LogError(ErrNoValidStorage.Error())
	}
}

// OnCommand implements the hook from the Bot
func (p *QuotesPlugin) OnCommand(cmd string, content string, post model.Post) {
	if post.IsPrivate {
		return
	}

	// TODO (#35): Use command prefix !quote for everything, e.g., !quote add some quote instead of !quoteadd
	if cmd == "quote" {
		p.onCommandQuote(content, post)
		return
	} else if cmd == "quoteadd" && len(content) > 0 {
		p.onCommandQuoteAdd(content, post)
		return
	} else if (cmd == "quoteadd" && len(content) == 0) || cmd == "quotehelp" {
		p.returnHelp(post.ChannelID)
		return
	}

	if !p.cfg.OnlyMods || utils.StringSliceContains(p.cfg.Mods, post.User.Name) {
		if cmd == "quoteremove" && len(content) > 0 {
			p.onCommandQuoteRemove(post)
			return
		} else if cmd == "quoteremove" {
			p.returnHelpRemove(post.ChannelID)
			return
		}
	} else {
		p.API.LogDebug("Not parsing as command, because User " + post.User.Name + " is not part of mods")
	}
}
