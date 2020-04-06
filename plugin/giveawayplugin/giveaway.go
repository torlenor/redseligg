package giveawayplugin

import (
	"math/rand"
	"time"

	"github.com/torlenor/abylebotter/plugin"
)

type randomizer interface {
	Intn(max int) int
	Shuffle(n int, swap func(i, j int))
}

// GiveawayPlugin is a plugin that lets you hold giveaways in your channel and let the bot pick a winner.
type GiveawayPlugin struct {
	plugin.AbyleBotterPlugin

	runningGiveaways runningGiveaways

	randomizer randomizer

	ticker         *time.Ticker
	tickerDoneChan chan bool
}

// New returns a new GiveawayPlugin
func New() (GiveawayPlugin, error) {
	ep := GiveawayPlugin{
		runningGiveaways: make(map[string]*giveaway),
		randomizer:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	return ep, nil
}
