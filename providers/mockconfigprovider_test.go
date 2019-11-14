package providers

import (
	"fmt"

	"github.com/torlenor/abylebotter/config"
)

type mockConfigProvider struct {
	pluginsConfig config.PluginConfigs
}

func (m *mockConfigProvider) GetBotConfig(id string) (config.BotConfig, error) {
	switch id {
	case "mockSlackID":
		return config.BotConfig{
			Type:    "mockSlack",
			Plugins: m.pluginsConfig,
		}, nil
	case "mockMattermostID":
		return config.BotConfig{
			Type:    "mockMattermost",
			Plugins: m.pluginsConfig,
		}, nil
	case "mockSomeOtherPlatformID":
		return config.BotConfig{
			Type:    "mockSomeOtherPlatform",
			Plugins: m.pluginsConfig,
		}, nil
	default:
		return config.BotConfig{}, fmt.Errorf("Unknown bot id %s", id)
	}
}
