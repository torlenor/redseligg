package discord

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/torlenor/redseligg/botconfig"

	"github.com/gorilla/websocket"

	"github.com/torlenor/redseligg/logging"
	"github.com/torlenor/redseligg/platform"
	"github.com/torlenor/redseligg/plugin"
	"github.com/torlenor/redseligg/storage"
	"github.com/torlenor/redseligg/utils"
	"github.com/torlenor/redseligg/webclient"
)

var (
	log = logging.Get("DiscordBot")
)

// Used for injection in unit tests
var newHeartbeatSender = func(ws webSocketClient) *discordHeartBeatSender {
	return &discordHeartBeatSender{ws: ws}
}

type webSocketClient interface {
	Dial(wsURL string) error
	Close() error

	ReadMessage() (int, []byte, error)

	SendMessage(messageType int, data []byte) error
	SendJSONMessage(v interface{}) error
}

type api interface {
	Call(path string, method string, body string) (webclient.APIResponse, error)
}

// The Bot struct holds parameters related to the bot
type Bot struct {
	platform.BotImpl

	storage storage.Storage

	api api

	gatewayURL string
	ws         webSocketClient

	knownChannels     map[string]channelCreate
	token             string
	ownSnowflakeID    string
	currentSeqNumber  int
	heartBeatStopChan chan bool
	seqNumberChan     chan int

	plugins []plugin.Hooks

	wg sync.WaitGroup

	watchdog *utils.Watchdog

	guilds        map[string]guildCreate // map[ID]
	guildNameToID map[string]string

	sessionID string
}

// CreateDiscordBotWithAPI creates a new instance of a DiscordBot with the
// provided api
func CreateDiscordBotWithAPI(api api, storage storage.Storage, cfg botconfig.DiscordConfig, ws webSocketClient) (*Bot, error) {
	log.Info("DiscordBot is CREATING itself")

	b := Bot{
		BotImpl: platform.BotImpl{
			ProvidedFeatures: map[string]bool{
				platform.FeatureMessagePost:    true,
				platform.FeatureMessageUpdate:  true,
				platform.FeatureMessageDelete:  true,
				platform.FeatureReactionNotify: true,
			},
		},

		api:     api,
		storage: storage,

		token: cfg.Token,
		ws:    ws,

		watchdog: &utils.Watchdog{},
	}

	url, err := b.getGateway()
	if err != nil {
		return nil, fmt.Errorf("Error connecting to Discord servers: %s", err)
	}
	b.gatewayURL = url

	b.knownChannels = make(map[string]channelCreate)

	b.seqNumberChan = make(chan int)

	b.guilds = make(map[string]guildCreate)
	b.guildNameToID = make(map[string]string)

	return &b, nil
}

// CreateDiscordBot creates a new instance of a DiscordBot
func CreateDiscordBot(cfg botconfig.DiscordConfig, storage storage.Storage, ws webSocketClient) (*Bot, error) {
	api := webclient.New("https://discordapp.com/api", "Bot "+cfg.Token, "application/json")

	return CreateDiscordBotWithAPI(api, storage, cfg, ws)
}

func (b *Bot) startHeartbeatSender(heartbeatInterval time.Duration) {
	b.heartBeatStopChan = make(chan bool)

	go func() {
		b.wg.Add(1)
		heartBeat(heartbeatInterval, newHeartbeatSender(b.ws), b.heartBeatStopChan, b.seqNumberChan, b.onFail)
		defer b.wg.Done()
	}()
	b.watchdog.SetFailCallback(b.onFail).Start(2 * heartbeatInterval)
}

func (b *Bot) openGatewayConnection() error {
	b.ws.Dial(b.gatewayURL)

	// WS: 1 - 10 HELLO
	_, message, err := b.ws.ReadMessage()
	if err != nil {
		return fmt.Errorf("Error occurred during initial communication with Discord Gateway: %s", err)
	}

	var data event
	if err := json.Unmarshal(message, &data); err != nil {
		return fmt.Errorf("Error occurred during initial communication with Discord Gateway: Could not unmarshal event: %s", err)
	}

	if data.Op != 10 { // If not Hello from Discord Gateway
		return fmt.Errorf("Error occurred during initial communication with Discord Gateway: Did not receive a HELLO, but OP Code %d", data.Op)
	}

	// Perform IDENT/RESUME
	if b.sessionID == "" {
		err = sendIdent(b.token, b.ws)
	} else {
		err = sendResume(b.token, b.sessionID, b.currentSeqNumber, b.ws)
	}
	if err != nil {
		return fmt.Errorf("Error occurred during initial communication with Discord Gateway: Could not send IDENT/RESUME: %s", err)
	}

	var helloEvent hello
	if err := json.Unmarshal(data.RawData, &helloEvent); err != nil {
		return fmt.Errorf("Error occurred during initial communication with Discord Gateway: Could not unmarshal HELLO event: %s", err)
	}

	// Start sending Heartbeats
	b.startHeartbeatSender(time.Duration(helloEvent.HeartbeatInterval) * time.Millisecond)

	return nil
}

func (b *Bot) run() {
	var e error
	defer func() {
		if e != nil {
			log.Error(e)
			go b.onFail()
		}
	}()

	e = b.openGatewayConnection()
	if e != nil {
		return
	}

	// Go into event handling
	for {
		_, message, err := b.ws.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				log.Tracef("Connection to Discord Gateway closed normally: %s", err)
				break
			} else if websocket.IsCloseError(err, websocket.CloseGoingAway) {
				log.Debugf("Received GoingAway from Discord Gateway, attempting a reconnect")
				b.stopHeartBeatWatchdog()
				b.ws.Close()
				err := b.openGatewayConnection()
				if err != nil {
					e = fmt.Errorf("Could not reconnect to Discord Gateway: %s", err.Error())
					break
				}
				continue
			} else {
				e = fmt.Errorf("Unhandled error in Discord Gateway communication logic: %s", err)
				break
			}
		}

		var data event
		if err := json.Unmarshal(message, &data); err != nil {
			log.Warnf("Could not unmarshal event from Discord Gateway: %s", err)
			continue
		}

		if data.Op == 7 { // Reconnect: You must reconnect with a new session immediately.
			log.Debugf("Received request to reconnect")
			b.stopHeartBeatWatchdog()
			b.ws.Close()
			e = b.openGatewayConnection()
		} else if data.Op == 9 { // Invalid Session: The session has been invalidated. You should reconnect and identify/resume accordingly.
			log.Warn("Invalid Session received")
			var invalidSessionData invalidSession
			json.Unmarshal(message, &invalidSessionData)
			if !invalidSessionData.D {
				// It does not want us to resume, so we are resetting our sessionID
				b.sessionID = ""
			}
			b.stopHeartBeatWatchdog()
			b.ws.Close()
			e = b.openGatewayConnection()
		} else if data.Op == 11 { // Heartbeat ACK
			b.watchdog.Feed()
		} else if data.Op == 0 { // Regular events to dispatch to event handlers
			switch data.Type {
			// Events managing communication with Discord API
			case "READY":
				b.handleReady(data.RawData)
			case "RESUMED":
				// Not needed for now
			case "RECONNECT":
				b.ws.Close()
				e = b.openGatewayConnection()
			// Real events
			case "MESSAGE_CREATE":
				b.handleMessageCreate(data.RawData)
			case "GUILD_CREATE":
				b.handleGuildCreate(data.RawData)
			case "PRESENCE_UPDATE":
				b.handlePresenceUpdate(data.RawData)
			case "PRESENCE_REPLACE":
				b.handlePresenceReplace(data.RawData)
			case "TYPING_START":
				b.handleTypingStart(data.RawData)
			case "CHANNEL_CREATE":
				b.handleChannelCreate(data.RawData)
			case "MESSAGE_REACTION_ADD":
				b.handleMessageReactionAdd(data.RawData)
			case "MESSAGE_REACTION_REMOVE":
				b.handleMessageReactionRemove(data.RawData)
			case "MESSAGE_DELETE":
				b.handleMessageDelete(data.RawData)
			case "MESSAGE_UPDATE":
				b.handleMessageUpdate(data.RawData)
			case "CHANNEL_PINS_UPDATE":
				b.handleChannelPinsUpdate(data.RawData)
			case "GUILD_MEMBER_UPDATE":
				b.handleGuildMemberUpdate(data.RawData)
			case "PRESENCES_REPLACE":
				b.handlePresencesReplace(data.RawData)
			default:
				log.Warnln("Unhandled message:", string(message))
			}
			b.currentSeqNumber = int(data.Seq)
			b.seqNumberChan <- b.currentSeqNumber
		} else {
			log.Warnf("Unknown Op Code %d received, data: %v", data.Op, data)
		}
	}
}

// Start the Discord Bot
func (b *Bot) start() error {
	log.Infof("DiscordBot is STARTING (have %d plugin(s))", len(b.plugins))

	err := b.ws.Dial(b.gatewayURL)
	if err != nil {
		return fmt.Errorf("Could not dial Discord WebSocket, Discord Bot not operational: %s", err.Error())
	}

	go func() {
		b.wg.Add(1)
		b.run()
		defer b.wg.Done()
	}()

	for _, plugin := range b.plugins {
		plugin.OnRun()
	}
	log.Info("DiscordBot is RUNNING")

	return nil
}

// Run the Discord Bot (blocking)
func (b *Bot) Run(ctx context.Context) error {
	err := b.start()
	if err != nil {
		return err
	}

	<-ctx.Done()

	for _, plugin := range b.plugins {
		plugin.OnStop()
	}

	b.stop()

	return nil
}

func (b *Bot) stopHeartBeatWatchdog() {
	b.watchdog.Stop()
	b.heartBeatStopChan <- true
}

// Stop the Discord Bot
func (b *Bot) stop() {
	log.Infoln("DiscordBot is SHUTING DOWN")

	b.stopHeartBeatWatchdog()

	err := b.sendCloseToWebsocket()
	if err != nil {
		log.Errorln("Error when writing close message to ws:", err)
	}

	b.wg.Wait()

	b.ws.Close()

	log.Infoln("DiscordBot is SHUT DOWN")
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
		Platform: "Discord",
		Healthy:  true,
		Plugins:  []platform.PluginInfo{},
	}
}

func (b *Bot) sendCloseToWebsocket() error {
	if b.ws != nil {
		return b.ws.SendMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	}
	return nil
}

func (b *Bot) sendCloseRestartToWebsocket() error {
	if b.ws != nil {
		return b.ws.SendMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseServiceRestart, ""))
	}
	return nil
}

func (b *Bot) onFail() {
	log.Warn("Encountered an error, trying to restart the bot")

	log.Debug("Stopping heartbeat watchdog")
	b.stopHeartBeatWatchdog()

	log.Debug("Sending close to websocket")
	err := b.sendCloseToWebsocket()
	if err != nil {
		log.Errorf("Error when writing close message to ws: %s, still trying to recover", err)
	}

	log.Debug("Waiting for waitgroup to be done")
	b.wg.Wait()

	log.Debug("Closing websocket")
	b.ws.Close()

	log.Debug("Get new gateway address")
	url, err := b.getGateway()
	if err != nil {
		log.Errorf("Error connecting to Discord servers: %s", err)
		return
	}
	b.gatewayURL = url

	log.Debugf("Dialing gateway at %s", b.gatewayURL)
	err = b.ws.Dial(b.gatewayURL)
	if err != nil {
		log.Errorln("Could not dial Discord WebSocket, Discord Bot not operational:", err)
		return
	}

	log.Debug("Launching run goroutine")
	go func() {
		b.wg.Add(1)
		b.run()
		defer b.wg.Done()
	}()

	log.Info("Recovery attempt finished")
}
