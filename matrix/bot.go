package matrix

import (
	"time"

	"github.com/torlenor/abylebotter/botinterface"
	"github.com/torlenor/abylebotter/config"
	"github.com/torlenor/abylebotter/events"
	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/plugins"
)

var (
	log = logging.Get("MatrixBot")
)

// The Bot struct holds parameters related to the bot
type Bot struct {
	api api

	sendMessageChan chan events.SendMessage

	pollingDone chan bool

	pollingInterval time.Duration

	receivers map[plugins.Plugin]chan events.ReceiveMessage

	knownPlugins []plugins.Plugin

	// We are wasting a little bit of memory and keep maps in both directions
	knownRooms   map[string]string // mapping of Room to RoomID
	knownRoomIDs map[string]string // mapping of RoomID to Room

	nextBatch string // contains the next batch to fetch in sync
}

// GetReceiveMessageChannel returns the channel which is used to notify
// about received messages from the bot
func (b *Bot) GetReceiveMessageChannel(plugin plugins.Plugin) chan events.ReceiveMessage {
	b.receivers[plugin] = make(chan events.ReceiveMessage)
	return b.receivers[plugin]
}

// GetSendMessageChannel returns the channel which is used to
// send messages using the bot
func (b *Bot) GetSendMessageChannel() chan events.SendMessage {
	return b.sendMessageChan
}

// The createMatrixBotWithAPI creates a new instance of a DiscordBot using the api interface api
func createMatrixBotWithAPI(api api, username string, password string, token string) (*Bot, error) {
	log.Printf("MatrixBot is CREATING itself")
	b := Bot{api: api}

	if len(token) == 0 {
		err := b.api.login(username, password)
		if err != nil {
			return nil, err
		}
	} else {
		// just use the provided access token
		b.api.updateAuthToken(token)
	}

	b.pollingDone = make(chan bool)
	b.pollingInterval = 1000 * time.Millisecond

	b.sendMessageChan = make(chan events.SendMessage)

	b.knownRooms = make(map[string]string)
	b.knownRoomIDs = make(map[string]string)

	b.receivers = make(map[plugins.Plugin]chan events.ReceiveMessage)

	return &b, nil
}

// CreateMatrixBot creates a new instance of a DiscordBot
func CreateMatrixBot(cfg config.MatrixConfig) (*Bot, error) {
	api := &matrixAPI{server: cfg.Server, authToken: cfg.Token}
	return createMatrixBotWithAPI(api, cfg.Username, cfg.Password, cfg.Token)
}

// Status returns the current status of MatrixBot
func (b *Bot) Status() botinterface.BotStatus {
	return botinterface.BotStatus{Running: true}
}

// AddPlugin adds the given plugin to the current bot
func (b *Bot) AddPlugin(plugin plugins.Plugin) {
	plugin.ConnectChannels(b.GetReceiveMessageChannel(plugin), b.GetSendMessageChannel())
	b.knownPlugins = append(b.knownPlugins, plugin)
	log.Debugln("Connected plugin", plugin.GetName())
}

func (b *Bot) addKnownRoom(roomID string, room string) {
	log.Debugln("Added new known Room:", roomID, room)
	b.knownRoomIDs[roomID] = room
	b.knownRooms[room] = roomID
}

func (b *Bot) removeKnownRoom(roomID string, room string) {
	log.Debugln("Removed known Room:", roomID, room)
	delete(b.knownRoomIDs, roomID)
	delete(b.knownRooms, room)
}

func (b *Bot) removeKnownRoomFromID(roomID string) {
	log.Debugln("Removed known Room with ID:", roomID)
	delete(b.knownRooms, b.knownRoomIDs[roomID])
	delete(b.knownRoomIDs, roomID)
}
