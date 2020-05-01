package providers

import (
	"fmt"

	"git.abyle.org/redseligg/botorchestrator/botconfig"
)

type mockConfigProvider struct {
	pluginsConfig botconfig.PluginConfigs
}

func (m *mockConfigProvider) GetBotConfig(id string) (botconfig.BotConfig, error) {
	switch id {
	case "mockSlackID":
		return botconfig.BotConfig{
			Type:    "mockSlack",
			Enabled: true,
			Plugins: m.pluginsConfig,
		}, nil
	case "mockMattermostID":
		return botconfig.BotConfig{
			Type:    "mockMattermost",
			Plugins: m.pluginsConfig,
		}, nil
	case "mockSomeOtherPlatformID":
		return botconfig.BotConfig{
			Type:    "mockSomeOtherPlatform",
			Plugins: m.pluginsConfig,
		}, nil
	default:
		return botconfig.BotConfig{}, fmt.Errorf("Unknown bot id %s", id)
	}
}

func (m *mockConfigProvider) GetAllEnabledBotIDs() []string {
	return []string{"mockSlack"}
}
