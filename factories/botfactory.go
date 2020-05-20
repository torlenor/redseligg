package factories

import (
	"fmt"

	"github.com/torlenor/redseligg/botconfig"
	"github.com/torlenor/redseligg/commanddispatcher"

	"github.com/torlenor/redseligg/logging"
	"github.com/torlenor/redseligg/platform"
	"github.com/torlenor/redseligg/platform/discord"
	"github.com/torlenor/redseligg/platform/matrix"
	"github.com/torlenor/redseligg/platform/mattermost"
	"github.com/torlenor/redseligg/platform/slack"
	"github.com/torlenor/redseligg/platform/twitch"
	"github.com/torlenor/redseligg/ws"
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

	logBotFactory.Tracef("Creating CommandDispatcher for botID %s", p)
	commandDispatcher := commanddispatcher.New(config.GeneralConfig.CallPrefix)

	switch p {
	case "slack":
		slackCfg, err := config.AsSlackConfig()
		if err != nil {
			return nil, fmt.Errorf("Error creating Slack bot: %s", err)
		}

		bot, err = slack.CreateSlackBot(slackCfg, storage, ws.NewClient())
		if err != nil {
			return nil, fmt.Errorf("Error creating Slack bot: %s", err)
		}
	case "mattermost":
		mmCfg, err := config.AsMattermostConfig()
		if err != nil {
			return nil, fmt.Errorf("Error creating Mattermost bot: %s", err)
		}

		bot, err = mattermost.CreateMattermostBot(mmCfg)
		if err != nil {
			return nil, fmt.Errorf("Error creating Mattermost bot: %s", err)
		}
	case "discord":
		discordCfg, err := config.AsDiscordConfig()
		if err != nil {
			return nil, fmt.Errorf("Error creating Discord bot: %s", err)
		}

		bot, err = discord.CreateDiscordBot(discordCfg, storage, commandDispatcher, ws.NewClient())
		if err != nil {
			return nil, fmt.Errorf("Error creating Discord bot: %s", err)
		}
	case "matrix":
		matrixCfg, err := config.AsMatrixConfig()
		if err != nil {
			return nil, fmt.Errorf("Error creating Matrix bot: %s", err)
		}

		bot, err = matrix.CreateMatrixBot(matrixCfg)
		if err != nil {
			return nil, fmt.Errorf("Error creating Matrix bot: %s", err)
		}
	case "twitch":
		twitchCfg, err := config.AsTwitchConfig()
		if err != nil {
			return nil, fmt.Errorf("Error creating Twitch bot: %s", err)
		}

		bot, err = twitch.CreateTwitchBot(twitchCfg, storage, ws.NewClient())
		if err != nil {
			return nil, fmt.Errorf("Error creating Twitch bot: %s", err)
		}
	default:
		return nil, fmt.Errorf("Unknown platform %s", p)
	}
	return bot, nil
}
