package providers

import (
	"fmt"

	"git.abyle.org/redseligg/botorchestrator/botconfig"

	"github.com/torlenor/abylebotter/platform"
)

// MockBotFactory can be used to generate bots for specific platforms
type MockBotFactory struct {
	bot MockBot
}

// CreateBot creates a new bot for the given platform with the provided configuration
func (b *MockBotFactory) CreateBot(p string, config botconfig.BotConfig) (platform.Bot, error) {
	var bot platform.Bot

	switch p {
	case "mockSlack":
		bot = &b.bot
	case "mockMattermost":
		bot = &b.bot
	default:
		return nil, fmt.Errorf("Unknown platform %s", p)
	}
	return bot, nil
}
