package giveawayplugin

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/torlenor/redseligg/model"
)

func parseTimeStringToDuration(timeStr string) (time.Duration, error) {
	duration, err := time.ParseDuration(timeStr)
	if err != nil {
		return time.Duration(0), fmt.Errorf("Giveaway duration invalid")
	}

	return duration, nil
}

func (p *GiveawayPlugin) returnHelp(channelID string) {
	p.returnMessage(channelID, "Type 'gstart <time> <secretword> [winners] [prize]' to start a new giveaway.")
}

func (p *GiveawayPlugin) returnMessage(channelID, msg string) {
	post := model.Post{
		ChannelID: channelID,
		Content:   msg,
	}
	p.API.CreatePost(post)
}

func (p *GiveawayPlugin) onCommandGStart(content string, post model.Post) {
	args := strings.Split(content, " ")

	if len(args) < 2 {
		p.returnHelp(post.ChannelID)
		return
	}

	timeStr := args[0]
	word := args[1]

	duration, err := parseTimeStringToDuration(timeStr)
	if err != nil {
		p.returnHelp(post.ChannelID)
		return
	}

	winners := 1
	if len(args) > 2 {
		if val, err := strconv.Atoi(args[2]); err == nil {
			winners = val
		} else {
			p.returnHelp(post.ChannelID)
			return
		}
	}

	prize := []string{}
	if len(args) > 3 {
		for _, arg := range args[3:] {
			prize = append(prize, arg)
		}
	}
	prizeStr := strings.Join(prize, " ")

	p.giveawaysMutex.Lock()
	defer p.giveawaysMutex.Unlock()

	if _, ok := p.runningGiveaways[post.ChannelID]; ok {
		p.returnMessage(post.ChannelID, "Giveaway already running.")
		return
	}

	giveaway := newGiveaway(post.ChannelID, post.User.ID, word, duration, winners, prizeStr)

	p.runningGiveaways[post.ChannelID] = &giveaway

	p.runningGiveaways[post.ChannelID].start(time.Now())
	p.returnMessage(post.ChannelID, "Giveaway started! Type "+word+" to participate.")
}

func (p *GiveawayPlugin) onCommandGEnd(post model.Post) {
	if g, ok := p.runningGiveaways[post.ChannelID]; ok {
		p.endGiveaway(g)
		return
	}

	p.returnMessage(post.ChannelID, "No giveaway running. Use gstart command to start a new one.")
}

func (p *GiveawayPlugin) onCommandGReroll(post model.Post) {
	p.giveawaysMutex.Lock()
	defer p.giveawaysMutex.Unlock()

	if _, ok := p.runningGiveaways[post.ChannelID]; ok {
		p.returnMessage(post.ChannelID, "Cannot pick a new winner. There is currently a giveaway running in this channel.")
		return
	}

	if g, ok := p.endedGiveaways[post.ChannelID]; ok {
		if len(g.participants) == 0 {
			p.returnMessage(g.channelID, "Cannot pick a new winner. There were no participants to the previous giveaway.")
			return
		}

		winnerID := g.getParticipantIDs()[p.randomizer.Intn(len(g.participants))]
		winner, err := g.getParticipant(winnerID)
		if err != nil {
			p.API.LogError("Something went wrong in rerolling a winner: " + err.Error())
			p.returnMessage(g.channelID, "There was an error in picking a new winner. Please start a new giveway. Sorry for the inconvenience.")
			return
		}

		endMessage := "The new winner is <@" + winner.ID + ">."

		if len(g.prize) > 0 {
			endMessage += " You won '" + g.prize + "'."
		}

		endMessage += " Congratulations!"

		p.returnMessage(g.channelID, endMessage)

		return
	}

	p.returnMessage(post.ChannelID, "No previous giveaway in that channel. Use gstart command to start a new one.")
}
