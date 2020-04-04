package giveawayplugin

import (
	"math/rand"
	"strings"
	"time"
)

type participant struct {
	ID   string
	Name string
}

type giveaway struct {
	channelID    string
	duration     time.Duration
	word         string
	numOfWinners int
	prize        string

	startTime time.Time

	participants   map[string]participant // [participantID] all participants in that giveaway
	participantIDs []string               // all participantIDs that take part in that giveaway
}

func newGiveaway(channelID, word string, duration time.Duration, numOfWinners int, prize string) giveaway {
	return giveaway{
		channelID:    channelID,
		duration:     duration,
		word:         word,
		numOfWinners: numOfWinners,
		prize:        prize,

		participants: make(map[string]participant),
	}
}

func (g *giveaway) start(startTime time.Time) {
	g.startTime = startTime
}

func (g giveaway) isFinished(currentTime time.Time) bool {
	return currentTime.Sub(g.startTime) > g.duration
}

func (g *giveaway) addParticipant(ID, name string) {
	if _, ok := g.participants[ID]; !ok {
		g.participants[ID] = participant{
			ID:   ID,
			Name: name,
		}
		g.participantIDs = append(g.participantIDs, ID)
	}
}

type runningGiveaways map[string]*giveaway // [channel]

func (p *GiveawayPlugin) endGiveaway(giveaway *giveaway) {
	delete(p.runningGiveaways, giveaway.channelID)

	var winners []string
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))

	// TODO roll various winners
	if len(giveaway.participants) == 0 {
		p.returnMessage(giveaway.channelID, "Cannot pick a winner. There were no participants to the giveaway.")
		return
	}
	winner := giveaway.participantIDs[randomizer.Intn(len(giveaway.participants))]
	winners = append(winners, "<@"+giveaway.participants[winner].ID+">")

	endMessage := "The winner(s) is/are " + strings.Join(winners, ", ") + "."

	if len(giveaway.prize) > 0 {
		endMessage += " You won '" + giveaway.prize + "'"
	}

	endMessage += " Congratulations!"

	p.returnMessage(giveaway.channelID, endMessage)
}
