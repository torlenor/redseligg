package voteplugin

import (
	"fmt"
	"strings"

	"github.com/torlenor/abylebotter/model"
)

const (
	HELP_TEXT = "Type '!vote What is the best color? [Red, Green, Blue]' to start a new giveaway.\nYou can omit the custom options in the [...] to initiate a simple Yes/No vote."
)

func (p *VotePlugin) returnHelp(channelID string) {
	p.returnMessage(channelID, HELP_TEXT)
}

func (p *VotePlugin) returnMessage(channelID, msg string) {
	post := model.Post{
		ChannelID: channelID,
		Content:   msg,
	}
	p.API.CreatePost(post)
}

func (p *VotePlugin) postAndStartVote(vote *vote) {
	post := vote.getCurrentPost()

	msgID, err := p.API.CreatePost(post)
	if err != nil {
		p.API.LogError("Something went wrong in creating the Vote message: " + err.Error())
		p.returnMessage(post.ChannelID, "Sorry to inform you, but we failed to create the Vote! Please try again later.")
		return
	}

	vote.messageIdent = msgID.PostedMessageIdent
}

func (p *VotePlugin) updatePost(vote *vote) {
	post := vote.getCurrentPost()

	_, err := p.API.UpdatePost(vote.messageIdent, post)
	if err != nil {
		p.API.LogError("Something went wrong in updating the Vote message: " + err.Error())
		return
	}
}

// onCommandVoteStart starts a new vote with the settings extracted
// from the received !vote command.
// Note: The command requires a valid !vote command. This check
// shall be performed at post retrieval.
func (p *VotePlugin) onCommandVoteStart(post model.Post) {
	cont := strings.Split(post.Content, " ")
	args := cont[1:]

	// TODO parse options
	// if empty, add Yes/No defaults
	description := strings.Join(args, " ")
	options := []string{"Yes", "No"}

	p.votesMutex.Lock()
	defer p.votesMutex.Unlock()
	nVote := newVote(voteSettings{
		ChannelID: post.ChannelID,
		Text:      description,
		Options:   options,
	})

	p.postAndStartVote(&nVote)
	p.runningVotes[nVote.Settings.Text] = &nVote
}

func (p *VotePlugin) onCommandVoteEnd(post model.Post) {
	cont := strings.Split(post.Content, " ")
	args := cont[1:]

	p.votesMutex.Lock()
	defer p.votesMutex.Unlock()
	description := strings.Join(args, " ")
	fmt.Printf("\nDescription = %s\n", description)
	for k, v := range p.runningVotes {
		fmt.Printf("k = %s, v: Text = %s, Channel = %s\n", k, v.Settings.Text, v.Settings.ChannelID)
	}
	if v, ok := p.runningVotes[description]; ok {
		if v.messageIdent.Channel == post.ChannelID {
			v.end()
			p.updatePost(v)
			delete(p.runningVotes, description)
			return
		}
	}

	p.returnMessage(post.ChannelID, "No vote running with that description in this channel. Use the !vote command to start a new one.")
}
