package twitch

import (
	"context"
	"fmt"
	"sync"

	"gopkg.in/irc.v3"

	"git.abyle.org/redseligg/botorchestrator/botconfig"

	"github.com/gorilla/websocket"
	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/platform"
	"github.com/torlenor/abylebotter/plugin"
	"github.com/torlenor/abylebotter/storage"
)

var (
	log = logging.Get("TwitchBot")
)

type webSocketClient interface {
	Dial(wsURL string) error
	Close() error

	ReadMessage() (int, []byte, error)

	SendMessage(messageType int, data []byte) error
	SendJSONMessage(v interface{}) error
}

// The Bot struct holds parameters related to the bot
type Bot struct {
	storage storage.Storage
	plugins []plugin.Hooks

	cfg botconfig.TwitchConfig

	ws webSocketClient

	wg sync.WaitGroup
}

// CreateTwitchBot creates a new instance of a TwitchBot
func CreateTwitchBot(cfg botconfig.TwitchConfig, storage storage.Storage, ws webSocketClient) (*Bot, error) {
	log.Info("TwitchBot is CREATING itself")

	b := Bot{
		storage: storage,
		cfg:     cfg,

		ws: ws,
	}

	return &b, nil
}

func (b *Bot) openWebSocketConnection() error {
	err := b.ws.Dial("wss://irc-ws.chat.twitch.tv:443")
	if err != nil {
		return err
	}

	b.ws.SendMessage(websocket.TextMessage, []byte("CAP REQ :twitch.tv/commands"))
	b.ws.SendMessage(websocket.TextMessage, []byte("CAP REQ :twitch.tv/membership"))
	b.ws.SendMessage(websocket.TextMessage, []byte("CAP REQ :twitch.tv/tags"))
	b.ws.SendMessage(websocket.TextMessage, []byte("PASS "+b.cfg.Token))
	b.ws.SendMessage(websocket.TextMessage, []byte("NICK #"+b.cfg.Username))
	for _, channel := range b.cfg.Channels {
		log.Infof("Joining channel " + channel)
		b.ws.SendMessage(websocket.TextMessage, []byte("JOIN #"+channel))
	}
	b.ws.SendMessage(websocket.TextMessage, []byte("USER #"+b.cfg.Username))

	return nil
}

func (b *Bot) messageLoop() {
	var e error
	defer func() {
		if e != nil {
			log.Error(e)
			// go b.onFail()
		}
	}()

	for {
		_, message, err := b.ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Debugln("Connection to Twitch Chat closed normally: ", err)
				break
			} else if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				// log.Infof("Received GoingAway from Twitch Chat, attempting a reconnect")
				// b.stopHeartBeatWatchdog()
				// b.ws.Close()
				// err := b.openGatewayConnection()
				// if err != nil {
				// 	e = fmt.Errorf("Could not reconnect to Twitch Gateway: %s", err.Error())
				// 	break
				// }
				// continue

				break
			} else {
				e = fmt.Errorf("Unhandled error in Twitch Chat communication logic: %s", err)
				break
			}
		}

		ircMessage, err := irc.ParseMessage(string(message))
		if err != nil {
			log.Errorf("Error parsing irc message: %s", err)
			continue
		}

		switch ircMessage.Command {
		case "PING":
			ircMessage.Command = "PONG"
			b.ws.SendMessage(websocket.TextMessage, []byte(ircMessage.String()))
		case "PRIVMSG":
			if len(ircMessage.Params) > 1 {
				post := model.Post{
					ChannelID: ircMessage.Params[0],
					User:      model.User{Name: ircMessage.User, ID: ircMessage.User},
					Content:   ircMessage.Params[1],
				}
				for _, plugin := range b.plugins {
					plugin.OnPost(post)
				}
			} else {
				log.Warnf("Params not long enough")
			}
		case "USERSTATE":
			// Not needed
		default:
			log.Warnf("Unhandled IRC command from server: %s, full message: %s", ircMessage.Command, ircMessage)
		}
	}
}

// Run the Bot (blocking)
func (b *Bot) Run(ctx context.Context) error {
	b.openWebSocketConnection()
	// RUN SOMETHING

	go func() {
		b.wg.Add(1)
		b.messageLoop()
		defer b.wg.Done()
	}()

	for _, plugin := range b.plugins {
		plugin.OnRun()
	}

	<-ctx.Done()

	for _, plugin := range b.plugins {
		plugin.OnStop()
	}

	// STOP SOMETHING
	b.ws.Close()

	return nil
}

// AddPlugin takes as argument a plugin and
// adds it to the bot providing it with the API
func (b *Bot) AddPlugin(plugin platform.BotPlugin) {
	plugin.SetAPI(b)
	b.plugins = append(b.plugins, plugin)
}

// GetInfo returns information about the Bot
func (b *Bot) GetInfo() platform.BotInfo {
	return platform.BotInfo{
		BotID:    "",
		Platform: "Twitch",
		Healthy:  true,
		Plugins:  []platform.PluginInfo{},
	}
}
