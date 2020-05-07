package voteplugin

import (
	"sync"

	"github.com/torlenor/abylebotter/botconfig"
	"github.com/torlenor/abylebotter/platform"
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
		AbyleBotterPlugin: plugin.AbyleBotterPlugin{
			NeededFeatures: []string{
				platform.FeatureMessagePost,
				platform.FeatureMessageUpdate,
				platform.FeatureReactionNotify,
			},
			Type: PLUGIN_TYPE,
		},
		cfg:          cfg,
		runningVotes: make(runningVotes),
	}

	return &ep, nil
}
