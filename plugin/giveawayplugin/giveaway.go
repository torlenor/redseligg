package giveawayplugin

import (
	"time"

	"github.com/torlenor/abylebotter/plugin"
)

// GiveawayPlugin is a plugin that lets you hold giveaways in your channel and let the bot pick a winner.
type GiveawayPlugin struct {
	plugin.AbyleBotterPlugin

	runningGiveaways runningGiveaways

	ticker         *time.Ticker
	tickerDoneChan chan bool
}

// New returns a new GiveawayPlugin
func New() (GiveawayPlugin, error) {
	ep := GiveawayPlugin{}

	ep.runningGiveaways = make(map[string]*giveaway)

	return ep, nil
}
