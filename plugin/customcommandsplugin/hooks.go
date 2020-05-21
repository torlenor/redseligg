package customcommandsplugin

import (
	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/utils"
)

// OnRun is called when the platform is ready
func (p *CustomCommandsPlugin) OnRun() {
	p.storage = p.getStorage()
	if p.storage == nil {
		p.API.LogError(ErrNoValidStorage.Error())
		return
	}

	commands, err := p.getCommands()
	if err != nil {
		p.API.LogError("Error getting commands for initial registration: " + err.Error())
		return
	}
	for _, command := range commands.Commands {
		p.API.RegisterCommand(p, command.Command)
	}

	p.API.RegisterCommand(p, "customcommand")
}

// OnCommand implements the hook from the Bot
func (p *CustomCommandsPlugin) OnCommand(cmd string, content string, post model.Post) {
	if post.IsPrivate {
		return
	}

	if !p.cfg.OnlyMods || utils.StringSliceContains(p.cfg.Mods, post.User.Name) {
		if cmd == "customcommand" && len(content) > 0 {
			p.onCommand(content, post)
			return
		} else if cmd == "customcommand" {
			p.returnHelp(post.ChannelID)
			return
		}
	} else {
		p.API.LogDebug("Not parsing as command, because User " + post.User.Name + " is not part of mods")
		return
	}

	p.onCustomCommand(cmd, post)
}
