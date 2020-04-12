package voteplugin

import (
	"github.com/torlenor/abylebotter/model"
)

func (p *VotePlugin) returnHelp(channelID string) {
	p.returnMessage(channelID, "Type '!vote What is the best color? [Red, Green, Blue]' to start a new giveaway.\nYou can omit the custom options in the [...] to initiate a simple Yes/No vote.")
}

func (p *VotePlugin) returnMessage(channelID, msg string) {
	post := model.Post{
		ChannelID: channelID,
		Content:   msg,
	}
	p.API.CreatePost(post)
}

func (p *VotePlugin) onCommandVoteStart(post model.Post) {
	// cont := strings.Split(post.Content, " ")
	// args := cont[1:]

	// p.returnMessage(post.ChannelID, "Giveaway started! Type "+word+" to participate.")
}

func (p *VotePlugin) onCommandVoteEnd(post model.Post) {
	// if g, ok := p.runningGiveaways[post.ChannelID]; ok {
	// 	p.endGiveaway(g)
	// 	return
	// }

	p.returnMessage(post.ChannelID, "No vote running. Use !vote command to start a new one.")
}
