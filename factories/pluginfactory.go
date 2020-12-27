package factories

import (
	"fmt"

	"github.com/torlenor/redseligg/botconfig"

	"github.com/torlenor/redseligg/platform"
	"github.com/torlenor/redseligg/plugin/archiveplugin"
	"github.com/torlenor/redseligg/plugin/customcommandsplugin"
	"github.com/torlenor/redseligg/plugin/echoplugin"
	"github.com/torlenor/redseligg/plugin/giveawayplugin"
	"github.com/torlenor/redseligg/plugin/httppingplugin"
	"github.com/torlenor/redseligg/plugin/quotesplugin"
	"github.com/torlenor/redseligg/plugin/rollplugin"
	"github.com/torlenor/redseligg/plugin/rssplugin"
	"github.com/torlenor/redseligg/plugin/timedmessagesplugin"
	"github.com/torlenor/redseligg/plugin/versionplugin"
	"github.com/torlenor/redseligg/plugin/voteplugin"
)

// PluginFactory can be used to generate plugins
type PluginFactory struct {
}

// CreatePlugin creates a new plugin with the provided configuration
func (b *PluginFactory) CreatePlugin(botID, pluginID string, pluginConfig botconfig.PluginConfig) (platform.BotPlugin, error) {
	var p platform.BotPlugin

	switch pluginConfig.Type {
	case "archive":
		p = &archiveplugin.ArchivePlugin{}
	case "customcommands":
		rp, err := customcommandsplugin.New(pluginConfig)
		if err != nil {
			return nil, err
		}
		p = rp
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
	case "rss":
		rp, err := rssplugin.New(pluginConfig)
		if err != nil {
			return nil, err
		}
		p = rp
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
		return nil, fmt.Errorf("Unknown plugin type '%s'", pluginConfig.Type)
	}

	p.SetBotPluginID(botID, pluginID)

	return p, nil
}
