package twitch

import (
	"context"
	"fmt"
	"sync"

	"gopkg.in/irc.v3"

	"github.com/torlenor/redseligg/botconfig"

	"github.com/gorilla/websocket"
	"github.com/torlenor/redseligg/logging"
	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/platform"
	"github.com/torlenor/redseligg/plugin"
	"github.com/torlenor/redseligg/storage"
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
	platform.BotImpl

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
		BotImpl: platform.BotImpl{
			ProvidedFeatures: map[string]bool{
				platform.FeatureMessagePost: true,
			},
		},

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
			go b.onFail()
		}
	}()

	for {
		_, message, err := b.ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Debugln("Connection to Twitch Chat closed normally: ", err)
				break
			} else if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Infof("Received GoingAway from Twitch Chat, attempting a reconnect")
				b.ws.Close()
				err := b.openWebSocketConnection()
				if err != nil {
					e = fmt.Errorf("Could not reconnect to Twitch Chat WebSocket: %s", err.Error())
					break
				}
				continue
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
		case "JOIN":
			// Viewer joins the channel
		case "PART":
			// View leaves the channel
		case "001":
			// Welcome message
		case "CAP":
			// Capabilities ack
		case "353":
			// List of current viewers "/NAMES"
		default:
			log.Warnf("Unhandled IRC command from server: %s, full message: %s", ircMessage.Command, ircMessage)
		}
	}
}

// Run the Bot (blocking)
func (b *Bot) Run(ctx context.Context) error {
	b.openWebSocketConnection()

	go func() {
		b.wg.Add(1)
		b.messageLoop()
		defer b.wg.Done()
	}()

	for _, plugin := range b.plugins {
		plugin.OnRun()
	}

	<-ctx.Done()
	log.Infoln("TwitchBot is SHUTING DOWN")

	for _, plugin := range b.plugins {
		plugin.OnStop()
	}

	err := b.sendCloseToWebsocket()
	if err != nil {
		log.Errorln("Error when writing close message to ws:", err)
	}

	b.wg.Wait()

	b.ws.Close()

	log.Infoln("TwitchBot is SHUT DOWN")

	return nil
}

func (b *Bot) sendCloseToWebsocket() error {
	if b.ws != nil {
		return b.ws.SendMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	}
	return nil
}

// AddPlugin takes as argument a plugin and
// adds it to the bot providing it with the API
func (b *Bot) AddPlugin(plugin platform.BotPlugin) {
	err := plugin.SetAPI(b)
	if err != nil {
		log.Errorf("Could not add plugin %s: %s", plugin.PluginType(), err)
	} else {
		b.plugins = append(b.plugins, plugin)
	}
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

func (b *Bot) onFail() {
	log.Warn("Encountered an error, trying to restart the bot")

	err := b.sendCloseToWebsocket()
	if err != nil {
		log.Errorf("Error when writing close message to ws: %s, still trying to recover", err)
	}

	b.wg.Wait()

	b.ws.Close()

	err = b.openWebSocketConnection()
	if err != nil {
		log.Errorln("Could not open Twitch Chat WebSocket, Twitch Bot not operational:", err)
		return
	}

	go func() {
		b.wg.Add(1)
		b.messageLoop()
		defer b.wg.Done()
	}()

	log.Info("Recovery attempt finished")
}
