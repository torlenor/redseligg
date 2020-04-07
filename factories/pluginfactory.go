package factories

import (
	"fmt"

	"git.abyle.org/redseligg/botorchestrator/botconfig"

	"github.com/torlenor/abylebotter/platform"
	"github.com/torlenor/abylebotter/plugin/echoplugin"
	"github.com/torlenor/abylebotter/plugin/giveawayplugin"
	"github.com/torlenor/abylebotter/plugin/httppingplugin"
	"github.com/torlenor/abylebotter/plugin/rollplugin"
	"github.com/torlenor/abylebotter/plugin/versionplugin"
)

// PluginFactory can be used to generate plugins
type PluginFactory struct {
}

// CreatePlugin creates a new plugin with the provided configuration
func (b *PluginFactory) CreatePlugin(plugin string, pluginConfig botconfig.PluginConfig) (platform.BotPlugin, error) {
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
	case "version":
		p = &versionplugin.VersionPlugin{}
	default:
		return nil, fmt.Errorf("Unknown plugin type %s", pluginConfig.Type)
	}

	return p, nil
}
