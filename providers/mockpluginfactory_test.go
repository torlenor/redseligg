package providers

import (
	"fmt"

	"github.com/torlenor/abylebotter/botconfig"

	"github.com/torlenor/abylebotter/platform"
)

// MockPluginFactory can be used to generate plugins
type MockPluginFactory struct {
	plugin MockPlugin
}

// CreatePlugin creates a new plugin with the provided configuration
func (b *MockPluginFactory) CreatePlugin(botID, pluginID string, pluginConfig botconfig.PluginConfig) (platform.BotPlugin, error) {
	var p platform.BotPlugin

	switch pluginConfig.Type {
	case "mockEcho":
		p = &b.plugin
	case "mockRoll":
		p = &b.plugin
	default:
		return nil, fmt.Errorf("Unknown plugin type %s", pluginConfig.Type)
	}

	return p, nil
}
