package voteplugin

import (
	"fmt"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/utils"
)

// OnRun implements the hook from the Boot
func (p *VotePlugin) OnRun() {
	// TODO: Use only one vote command instead of three separate ones
	p.API.RegisterCommand(p, "vote")
	p.API.RegisterCommand(p, "voteend")
	p.API.RegisterCommand(p, "votehelp")
}

// OnCommand implements the hook from the Bot
func (p *VotePlugin) OnCommand(cmd string, content string, post model.Post) {
	if post.IsPrivate {
		return
	}

	if !p.cfg.OnlyMods || utils.StringSliceContains(p.cfg.Mods, post.User.Name) {
		if cmd == "vote" && len(content) > 0 {
			p.onCommandVoteStart(content, post)
			return
		} else if cmd == "voteend" && len(content) > 0 {
			p.onCommandVoteEnd(content, post)
			return
		} else if cmd == "voteend" {
			p.returnVoteEndHelp(post.ChannelID)
			return
		} else if (cmd == "vote" && len(content) == 0) || cmd == "votehelp" {
			p.returnHelp(post.ChannelID)
			return
		}
	} else {
		p.API.LogDebug("Not parsing as command, because User " + post.User.Name + " is not part of mods")
	}
}

// if not found returns nil
func (p *VotePlugin) getVoteForMessageIdent(messageIdent model.MessageIdentifier) *vote {
	if k, ok := p.runningVotes[messageIdent.Channel]; ok {
		for _, v := range k {
			if v.messageIdent.Channel == messageIdent.Channel && v.messageIdent.ID == messageIdent.ID {
				return v
			}
		}
	}

	return nil
}

// OnReactionAdded implements the hook from the bot
func (p *VotePlugin) OnReactionAdded(reaction model.Reaction) {
	p.API.LogDebug(fmt.Sprintf("Received ReactionAdded: %v", reaction))

	if v := p.getVoteForMessageIdent(reaction.Message); v != nil {
		if v.countVote(reaction.Reaction) {
			p.updatePost(v)
		}
	}
}

// OnReactionRemoved implements the hook from the bot
func (p *VotePlugin) OnReactionRemoved(reaction model.Reaction) {
	p.API.LogDebug(fmt.Sprintf("Received ReactionRemoved: %v", reaction))

	if v := p.getVoteForMessageIdent(reaction.Message); v != nil {
		if v.removeVote(reaction.Reaction) {
			p.updatePost(v)
		}
	}
}
