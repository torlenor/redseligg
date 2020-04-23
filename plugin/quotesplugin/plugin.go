package quotesplugin

import (
	"math/rand"
	"time"

	"git.abyle.org/redseligg/botorchestrator/botconfig"
	"github.com/torlenor/abylebotter/plugin"
)

const (
	PLUGIN_TYPE = "quotes"
)

type randomizer interface {
	Intn(max int) int
}

// QuotesPlugin is a plugin that allows viewers or mods to add quotes and randomly fetch one.
type QuotesPlugin struct {
	plugin.AbyleBotterPlugin

	cfg config

	randomizer randomizer
}

// New returns a new QuotesPlugin
func New(pluginConfig botconfig.PluginConfig) (*QuotesPlugin, error) {
	cfg, err := parseConfig(pluginConfig)
	if err != nil {
		return nil, err
	}

	ep := QuotesPlugin{
		cfg:        cfg,
		randomizer: rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	return &ep, nil
}
