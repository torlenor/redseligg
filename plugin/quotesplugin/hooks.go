package quotesplugin

import (
	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/utils"
)

// OnRun is called when the platform is ready
func (p *QuotesPlugin) OnRun() {
	p.API.RegisterCommand(p, command)

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

	subcommand, argument := utils.ExtractSubCommandAndArgsString(content)

	if subcommand != "add" && subcommand != "remove" && subcommand != "help" {
		// In case of raw quote command the subcommand is the actual argument to fetch
		// a specific quote
		p.onCommandQuote(subcommand, post)
		return
	} else if subcommand == "add" && len(argument) > 0 {
		p.onCommandQuoteAdd(argument, post)
		return
	} else if (subcommand == "add" && len(argument) == 0) || subcommand == "help" {
		p.returnHelp(post.ChannelID)
		return
	}

	if !p.cfg.OnlyMods || utils.StringSliceContains(p.cfg.Mods, post.User.Name) {
		if subcommand == "remove" && len(argument) > 0 {
			p.onCommandQuoteRemove(argument, post)
			return
		} else if subcommand == "remove" {
			p.returnHelpRemove(post.ChannelID)
			return
		}
	} else {
		p.API.LogDebug("Not parsing as command, because User " + post.User.Name + " is not part of mods")
	}
}
