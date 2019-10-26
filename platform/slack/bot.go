package slack

import (
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/torlenor/abylebotter/botinterface"
	"github.com/torlenor/abylebotter/config"
	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/platform"
	"github.com/torlenor/abylebotter/plugin"
	"github.com/torlenor/abylebotter/utils"
)

type webSocketClient interface {
	Dial(wsURL string) error
	Stop()

	ReadMessage() (int, []byte, error)

	SendMessage(messageType int, data []byte) error
	SendJSONMessage(v interface{}) error
}

// The Bot struct holds parameters related to the bot
type Bot struct {
	config config.SlackConfig
	log    *logrus.Entry

	rtmURL string
	ws     webSocketClient

	channels channelManager
	users    userManager

	plugins []plugin.Hooks

	wg sync.WaitGroup

	pingSenderStop chan bool

	watchdog *watchdog

	idProvider utils.IDProvider
}

// CreateSlackBot creates a new instance of a SlackBot
func CreateSlackBot(cfg config.SlackConfig, ws webSocketClient) (*Bot, error) {
	log := logging.Get("SlackBot")
	log.Printf("SlackBot is CREATING itself")

	b := Bot{
		config: cfg,
		log:    log,

		ws: ws,

		channels: newChannelManager(),
		users:    newUserManager(),

		watchdog: &watchdog{},
	}

	if len(b.config.Token) == 0 {
		return nil, fmt.Errorf("No Slack token defined in config file")
	}

	rtmConnectResponse, err := b.RtmConnect()
	if err != nil {
		return nil, fmt.Errorf("Error connecting to Slack servers: %s", err)
	}

	b.rtmURL = rtmConnectResponse.URL

	return &b, nil
}

func (b *Bot) startPingWatchdog() {
	b.pingSenderStop = make(chan bool)
	go func() {
		b.wg.Add(1)
		pingSender(5*time.Second, b.sendPing, b.pingSenderStop)
		defer b.wg.Done()
	}()
	b.watchdog.SetFailCallback(b.onFail).Start(10 * time.Second)
}

// Start the Bot
func (b *Bot) Start() {
	b.log.Infof("SlackBot is STARTING (have %d plugin(s))", len(b.plugins))

	err := b.ws.Dial(b.rtmURL)
	if err != nil {
		b.log.Errorln("Could not dial Slack RTM WebSocket, Slack Bot not operational:", err)
		return
	}

	err = b.populateChannelList()
	if err != nil {
		b.log.Warnln("Populating Channel List failed, no Channel information will be available:", err)
	}

	err = b.populateUserList()
	if err != nil {
		b.log.Warnln("Populating User List failed, no User information will be available:", err)
	}

	b.startPingWatchdog()

	go func() {
		b.wg.Add(1)
		b.run()
		defer b.wg.Done()
	}()

	b.log.Infoln("SlackBot is RUNNING")
}

func (b *Bot) stopPingWatchdog() {
	b.watchdog.Stop()
	b.pingSenderStop <- true
}

// Stop the Bot
func (b *Bot) Stop() {
	b.log.Infoln("SlackBot is SHUTING DOWN")

	b.stopPingWatchdog()

	err := b.ws.SendMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		b.log.Warnln("Error when writing close message to ws:", err)
	}

	b.wg.Wait()
	b.ws.Stop()

	b.log.Infoln("SlackBot is SHUT DOWN")
}

// Status returns the current status of the SlackBot
func (b *Bot) Status() botinterface.BotStatus {
	status := botinterface.BotStatus{
		Running: true,
		Fail:    false,
		Fatal:   false,
	}
	return status
}

// AddPlugin takes as argument a plugin and
// adds it to the bot providing it with the API
func (b *Bot) AddPlugin(plugin platform.BotPlugin) {
	plugin.SetAPI(b)
	b.plugins = append(b.plugins, plugin)
}