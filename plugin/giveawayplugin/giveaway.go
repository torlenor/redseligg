package giveawayplugin

import (
	"math/rand"
	"sync"
	"time"

	"git.abyle.org/redseligg/botorchestrator/botconfig"

	"github.com/torlenor/abylebotter/plugin"
)

type randomizer interface {
	Intn(max int) int
	Shuffle(n int, swap func(i, j int))
}

// GiveawayPlugin is a plugin that lets you hold giveaways in your channel and let the bot pick a winner.
type GiveawayPlugin struct {
	plugin.AbyleBotterPlugin

	cfg config

	giveawaysMutex   sync.Mutex
	runningGiveaways runningGiveaways
	endedGiveaways   runningGiveaways

	randomizer randomizer

	ticker         *time.Ticker
	tickerDoneChan chan bool
}

// New returns a new GiveawayPlugin
func New(pluginConfig botconfig.PluginConfig) (*GiveawayPlugin, error) {
	cfg, err := parseConfig(pluginConfig)
	if err != nil {
		return nil, err
	}

	ep := GiveawayPlugin{
		cfg:              cfg,
		runningGiveaways: make(map[string]*giveaway),
		endedGiveaways:   make(map[string]*giveaway),
		randomizer:       rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	return &ep, nil
}
