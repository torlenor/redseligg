package matrix

import (
	"botinterface"
	"logging"
	"plugins"
	"time"

	"events"
)

var (
	log = logging.Get("MatrixBot")
)

// The Bot struct holds parameters related to the bot
type Bot struct {
	api api

	receiveMessageChan chan events.ReceiveMessage
	sendMessageChan    chan events.SendMessage
	commandChan        chan events.Command
	pollingDone        chan bool

	pollingInterval time.Duration

	knownPlugins []plugins.Plugin

	// We are wasting a little bit of memory and keep maps in both directions
	knownRooms   map[string]string // mapping of Room to RoomID
	knownRoomIDs map[string]string // mapping of RoomID to Room

	nextBatch string // contains the next batch to fetch in sync
}

// GetReceiveMessageChannel returns the channel which is used to notify
// about received messages from the bot. For DiscordBot these messages
// can be normal channel messages, whispers
func (b Bot) GetReceiveMessageChannel() chan events.ReceiveMessage {
	return b.receiveMessageChan
}

// GetSendMessageChannel returns the channel which is used to
// send messages using the bot. For DiscordBot these messages
// can be normal channel messages, whispers
func (b Bot) GetSendMessageChannel() chan events.SendMessage {
	return b.sendMessageChan
}

// GetCommandChannel gives a channel to control the bot from
// a plugin
func (b Bot) GetCommandChannel() chan events.Command {
	return b.commandChan
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

	b.api.connectToMatrixServer()

	b.pollingDone = make(chan bool)
	b.pollingInterval = 1000 * time.Millisecond

	b.receiveMessageChan = make(chan events.ReceiveMessage)
	b.sendMessageChan = make(chan events.SendMessage)
	b.commandChan = make(chan events.Command)

	b.knownRooms = make(map[string]string)
	b.knownRoomIDs = make(map[string]string)

	return &b, nil
}

// CreateMatrixBot creates a new instance of a DiscordBot
func CreateMatrixBot(server string, username string, password string, token string) (*Bot, error) {
	api := &matrixAPI{server: server, authToken: token}
	return createMatrixBotWithAPI(api, username, password, token)
}

// Status returns the current status of MatrixBot
func (b *Bot) Status() botinterface.BotStatus {
	return botinterface.BotStatus{Running: true}
}

// AddPlugin adds the give plugin to the current bot
func (b *Bot) AddPlugin(plugin plugins.Plugin) {
	plugin.ConnectChannels(b.GetReceiveMessageChannel(), b.GetSendMessageChannel(), b.GetCommandChannel())
	b.knownPlugins = append(b.knownPlugins, plugin)
}

func (b *Bot) updateRoom(roomID string, room string) {
	b.knownRooms[room] = roomID
	b.knownRoomIDs[roomID] = room
}
