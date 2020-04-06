package giveawayplugin

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

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
	initiatorID  string

	startTime time.Time

	participantsMutex sync.Mutex
	participants      map[string]participant // [participantID] all participants in that giveaway
	participantIDs    []string               // all participantIDs that take part in that giveaway
}

func newGiveaway(channelID, initiatorID, word string, duration time.Duration, numOfWinners int, prize string) giveaway {
	return giveaway{
		channelID:    channelID,
		duration:     duration,
		word:         word,
		numOfWinners: numOfWinners,
		prize:        prize,

		initiatorID: initiatorID,

		participants: make(map[string]participant),
	}
}

func (g *giveaway) start(startTime time.Time) {
	g.startTime = startTime
}

func (g *giveaway) isFinished(currentTime time.Time) bool {
	return currentTime.Sub(g.startTime) > g.duration
}

func (g *giveaway) addParticipant(ID, name string) {
	g.participantsMutex.Lock()
	defer g.participantsMutex.Unlock()

	if _, ok := g.participants[ID]; !ok {
		g.participants[ID] = participant{
			ID:   ID,
			Name: name,
		}
		g.participantIDs = append(g.participantIDs, ID)
	}
}

func (g *giveaway) getParticipantIDs() []string {
	g.participantsMutex.Lock()
	defer g.participantsMutex.Unlock()

	cpy := make([]string, len(g.participantIDs))
	copy(cpy, g.participantIDs)
	return cpy
}
func (g *giveaway) getParticipant(participantID string) (participant, error) {
	g.participantsMutex.Lock()
	defer g.participantsMutex.Unlock()

	if p, ok := g.participants[participantID]; ok {
		return p, nil
	}
	return participant{}, fmt.Errorf("Participant not found")
}

type runningGiveaways map[string]*giveaway // [channel]

func (p *GiveawayPlugin) endGiveaway(giveaway *giveaway) {
	delete(p.runningGiveaways, giveaway.channelID)

	var winners []string

	if len(giveaway.participants) == 0 {
		p.returnMessage(giveaway.channelID, "Cannot pick a winner. There were no participants to the giveaway.")
		return
	}

	participants := giveaway.getParticipantIDs()
	p.randomizer.Shuffle(len(participants), func(i, j int) { participants[i], participants[j] = participants[j], participants[i] })
	for i := 0; i < min(giveaway.numOfWinners, len(participants)); i++ {
		winner, err := giveaway.getParticipant(participants[i])
		if err != nil {
			p.API.LogError("Something went wrong in picking a winner: " + err.Error())
			continue
		}
		winners = append(winners, "<@"+winner.ID+">")
	}

	endMessage := "The winner(s) is/are " + strings.Join(winners, ", ") + "."

	if len(giveaway.prize) > 0 {
		endMessage += " You won '" + giveaway.prize + "'."
	}

	endMessage += " Congratulations!"

	p.returnMessage(giveaway.channelID, endMessage)
}
