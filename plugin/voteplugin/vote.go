package voteplugin

import (
	"fmt"
	"strconv"

	"github.com/torlenor/redseligg/model"
)

// TODO: There has to be a mapping from that to the emojis and syntax used from the platform via an API call.
func getDefaultReactions() []string {
	return []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "keycap_ten"}
}

type runningVotes map[string]map[string]*vote // [Channel][Description]

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

	Ended bool

	// messageIdent stores the info about the created message
	// and will be used to update the message with the current
	// or end results
	messageIdent model.MessageIdentifier
}

func newVote(voteSettings voteSettings) (vote, error) {
	options := []*voteOption{}

	defaultReactions := getDefaultReactions()

	if len(voteSettings.Options) > len(defaultReactions) {
		return vote{}, fmt.Errorf("More than the allowed number of options specified. Please specify " + strconv.Itoa(len(defaultReactions)) + " or less options.")
	}

	for i, option := range voteSettings.Options {
		options = append(options, &voteOption{Description: option, AssignedReaction: defaultReactions[i]})
	}

	return vote{
		Settings: voteSettings,
		Options:  options,
	}, nil
}

// end ends a vote and posting the final results.
func (v *vote) end() {
	v.Ended = true
}

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

	if !v.Ended {
		content += "Participate by reacting with the appropriate emoji corresponding to the option you want to vote for!"
	} else {
		content += "This vote has ended, thanks for participating!"
	}

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
	if v.Ended {
		return false
	}

	for _, option := range v.Options {
		if option.AssignedReaction == reaction {
			option.Votes++
			return true
		}
	}
	return false
}

func (v *vote) removeVote(reaction string) bool {
	if v.Ended {
		return false
	}

	for _, option := range v.Options {
		if option.AssignedReaction == reaction {
			option.Votes--
			return true
		}
	}
	return false
}
