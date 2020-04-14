package voteplugin

import (
	"sync"

	"git.abyle.org/redseligg/botorchestrator/botconfig"
	"github.com/torlenor/abylebotter/plugin"
)

const (
	PLUGIN_TYPE = "vote"
)

// VotePlugin is a plugin to initiate a vote in the channel about arbitrary topics.
type VotePlugin struct {
	plugin.AbyleBotterPlugin

	cfg config

	votesMutex   sync.Mutex
	runningVotes runningVotes
}

// New returns a new VotePlugin
func New(pluginConfig botconfig.PluginConfig) (*VotePlugin, error) {
	cfg, err := parseConfig(pluginConfig)
	if err != nil {
		return nil, err
	}

	ep := VotePlugin{
		cfg:          cfg,
		runningVotes: make(runningVotes),
	}

	return &ep, nil
}
