package customcommandsplugin

import (
	"strings"

	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/utils"
)

// OnRun is called when the platform is ready
func (p *CustomCommandsPlugin) OnRun() {
	p.storage = p.getStorage()
	if p.storage == nil {
		p.API.LogError(ErrNoValidStorage.Error())
		return
	}
}

// OnPost implements the hook from the Bot
func (p *CustomCommandsPlugin) OnPost(post model.Post) {
	if post.IsPrivate {
		return
	}

	msg := strings.Trim(post.Content, " ")

	if !p.cfg.OnlyMods || utils.StringSliceContains(p.cfg.Mods, post.User.Name) {
		if strings.HasPrefix(msg, "!customcommand ") {
			p.onCommand(post)
			return
		} else if msg == "!customcommand" {
			p.returnHelp(post.ChannelID)
			return
		}
	} else {
		p.API.LogDebug("Not parsing as command, because User " + post.User.Name + " is not part of mods")
	}

	if strings.HasPrefix(msg, "!") {
		p.onCustomCommand(post)
	}
}
