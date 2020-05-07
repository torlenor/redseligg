package factories

import (
	"fmt"

	"github.com/torlenor/abylebotter/botconfig"

	"github.com/torlenor/abylebotter/platform"
	"github.com/torlenor/abylebotter/plugin/echoplugin"
	"github.com/torlenor/abylebotter/plugin/giveawayplugin"
	"github.com/torlenor/abylebotter/plugin/httppingplugin"
	"github.com/torlenor/abylebotter/plugin/quotesplugin"
	"github.com/torlenor/abylebotter/plugin/rollplugin"
	"github.com/torlenor/abylebotter/plugin/timedmessagesplugin"
	"github.com/torlenor/abylebotter/plugin/versionplugin"
	"github.com/torlenor/abylebotter/plugin/voteplugin"
)

// PluginFactory can be used to generate plugins
type PluginFactory struct {
}

// CreatePlugin creates a new plugin with the provided configuration
func (b *PluginFactory) CreatePlugin(botID, pluginID string, pluginConfig botconfig.PluginConfig) (platform.BotPlugin, error) {
	var p platform.BotPlugin

	switch pluginConfig.Type {
	case "echo":
		p = &echoplugin.EchoPlugin{}
	case "giveaway":
		rp, err := giveawayplugin.New(pluginConfig)
		if err != nil {
			return nil, err
		}
		p = rp
	case "roll":
		rp, err := rollplugin.New()
		if err != nil {
			return nil, err
		}
		p = &rp
	case "httpping":
		p = &httppingplugin.HTTPPingPlugin{}
	case "quotes":
		rp, err := quotesplugin.New(pluginConfig)
		if err != nil {
			return nil, err
		}
		p = rp
	case "timedmessages":
		rp, err := timedmessagesplugin.New(pluginConfig)
		if err != nil {
			return nil, err
		}
		p = rp
	case "version":
		p = &versionplugin.VersionPlugin{}
	case "vote":
		rp, err := voteplugin.New(pluginConfig)
		if err != nil {
			return nil, err
		}
		p = rp
	default:
		return nil, fmt.Errorf("Unknown plugin type %s", pluginConfig.Type)
	}

	p.SetBotPluginID(botID, pluginID)

	return p, nil
}
