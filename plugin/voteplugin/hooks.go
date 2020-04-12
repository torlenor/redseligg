package voteplugin

import (
	"strings"

	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/utils"
)

// OnPost implements the hook from the Bot
func (p *VotePlugin) OnPost(post model.Post) {
	if post.IsPrivate {
		return
	}

	msg := strings.Trim(post.Content, " ")
	if !p.cfg.OnlyMods || utils.StringSliceContains(p.cfg.Mods, post.User.Name) {
		if strings.HasPrefix(msg, "!vote ") {
			p.onCommandVoteStart(post)
			return
		} else if msg == "!voteend" {
			p.onCommandVoteEnd(post)
			return
		} else if msg == "!vote" || msg == "!votehelp" {
			p.returnHelp(post.ChannelID)
			return
		}
	} else {
		p.API.LogDebug("Not parsing as command, because User " + post.User.Name + " is not part of mods")
	}
}
