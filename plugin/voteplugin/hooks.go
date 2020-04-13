package voteplugin

import (
	"fmt"
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
		} else if strings.HasPrefix(msg, "!voteend ") {
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

// if not found returns nil
func (p *VotePlugin) getVoteForMessageIdent(messageIdent model.MessageIdentifier) *vote {
	for _, v := range p.runningVotes {
		if v.messageIdent.Channel == messageIdent.Channel && v.messageIdent.ID == messageIdent.ID {
			return v
		}
	}

	fmt.Printf("Not found: %v\n", messageIdent)

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
