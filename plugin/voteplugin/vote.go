package voteplugin

import (
	"strconv"

	"github.com/torlenor/abylebotter/model"
)

// TODO: There has to be a mapping from that to the emojis and syntax used from the platform via an API call.
func getDefaultReactions() []string {
	return []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "then"}
}

type runningVotes map[string]*vote // [Description]

type voteSettings struct {
	ChannelID string   // ChannelID where the vote shall be started
	Text      string   // Text is the text that should be presented alongside that vote
	Options   []string // Options is the list of options requested for that vote
}

type voteOption struct {
	Description      string // Description of the voting option
	AssignedReaction string // The reaction assigned to that option
	Votes            uint32 // Number of votes received
}

type vote struct {
	Settings voteSettings
	Options  []*voteOption

	// messageIdent stores the info about the created message
	// and will be used to update the message with the current
	// or end results
	messageIdent model.MessageIdentifier
}

func newVote(voteSettings voteSettings) vote {
	options := []*voteOption{}

	defaultReactions := getDefaultReactions()

	for i, option := range voteSettings.Options {
		options = append(options, &voteOption{Description: option, AssignedReaction: defaultReactions[i]})
	}

	return vote{
		Settings: voteSettings,
		Options:  options,
	}
}

// start is initiating the voting and starts
// updating the post when there are new votes.
func (v *vote) start() {}

// end ends a vote and posting the final results.
func (v *vote) end() {}

func (v *vote) getContent() string {
	var content string

	content = "\n" + "*" + v.Settings.Text + "*"
	content += "\n"

	for _, option := range v.Options {
		content += ":" + option.AssignedReaction + ":" + ": " + option.Description
		if option.Votes != 0 {
			content += " - " + strconv.Itoa(int(option.Votes))
		}
		content += "\n"
	}

	content += "Participate by reacting with the appropriate emoji corresponding to the option you want to vote for!"

	return content
}

// getCurrentPost returns the post for the vote.
// On first call it should be used in createPost and
// after the vote is started UpdatePost shall be used
// to post it so that the previous post is replaced.
func (v *vote) getCurrentPost() model.Post {
	return model.Post{
		ChannelID: v.Settings.ChannelID,
		Content:   v.getContent(),
	}
}

func (v *vote) countVote(reaction string) bool {
	for _, option := range v.Options {
		if option.AssignedReaction == reaction {
			option.Votes++
			return true
		}
	}
	return false
}

func (v *vote) removeVote(reaction string) bool {
	for _, option := range v.Options {
		if option.AssignedReaction == reaction {
			option.Votes--
			return true
		}
	}
	return false
}
