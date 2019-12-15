package factories

import (
	"fmt"

	"git.abyle.org/reseligg/botorchestrator/botconfig"

	"github.com/torlenor/abylebotter/platform"
	"github.com/torlenor/abylebotter/plugin/echoplugin"
	"github.com/torlenor/abylebotter/plugin/httppingplugin"
	"github.com/torlenor/abylebotter/plugin/rollplugin"
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
	case "roll":
		rp, err := rollplugin.New()
		if err != nil {
			return nil, err
		}
		p = &rp
	case "httpping":
		p = &httppingplugin.HTTPPingPlugin{}
	default:
		return nil, fmt.Errorf("Unknown plugin type %s", pluginConfig.Type)
	}

	return p, nil
}
