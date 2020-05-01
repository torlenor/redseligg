package factories

import (
	"fmt"

	"git.abyle.org/redseligg/botorchestrator/botconfig"

	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/platform"
	"github.com/torlenor/abylebotter/platform/discord"
	"github.com/torlenor/abylebotter/platform/matrix"
	"github.com/torlenor/abylebotter/platform/mattermost"
	"github.com/torlenor/abylebotter/platform/slack"
	"github.com/torlenor/abylebotter/ws"
)

var (
	logBotFactory = logging.Get("BotFactory")
)

// BotFactory can be used to generate bots for specific platforms
type BotFactory struct{}

// CreateBot creates a new bot for the given platform with the provided configuration
func (b *BotFactory) CreateBot(p string, config botconfig.BotConfig) (platform.Bot, error) {
	var bot platform.Bot

	storageFactory := StorageFactory{}

	logBotFactory.Tracef("Creating storage backend for botID %s", p)
	storage, err := storageFactory.CreateBackend(config.StorageConfig)
	if err != nil {
		return nil, fmt.Errorf("Error creating storage backend for botID %s: %s", p, err)
	}

	switch p {
	case "slack":
		slackCfg, err := config.AsSlackConfig()
		if err != nil {
			return nil, fmt.Errorf("Error creating slack bot: %s", err)
		}

		bot, err = slack.CreateSlackBot(slackCfg, ws.NewClient())
		if err != nil {
			return nil, fmt.Errorf("Error creating slack bot: %s", err)
		}
	case "mattermost":
		mmCfg, err := config.AsMattermostConfig()
		if err != nil {
			return nil, fmt.Errorf("Error creating mattermost bot: %s", err)
		}

		bot, err = mattermost.CreateMattermostBot(mmCfg)
		if err != nil {
			return nil, fmt.Errorf("Error creating mattermost bot: %s", err)
		}
	case "discord":
		discordCfg, err := config.AsDiscordConfig()
		if err != nil {
			return nil, fmt.Errorf("Error creating discord bot: %s", err)
		}

		bot, err = discord.CreateDiscordBot(discordCfg, storage, ws.NewClient())
		if err != nil {
			return nil, fmt.Errorf("Error creating discord bot: %s", err)
		}
	case "matrix":
		matrixCfg, err := config.AsMatrixConfig()
		if err != nil {
			return nil, fmt.Errorf("Error creating discord bot: %s", err)
		}

		bot, err = matrix.CreateMatrixBot(matrixCfg)
		if err != nil {
			return nil, fmt.Errorf("Error creating discord bot: %s", err)
		}
	default:
		return nil, fmt.Errorf("Unknown platform %s", p)
	}
	return bot, nil
}
