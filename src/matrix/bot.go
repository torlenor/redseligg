package matrix

import (
	"botinterface"
	"io/ioutil"
	"logging"
	"net/http"
	"plugins"
	"strings"
	"time"

	"events"
)

var (
	log = logging.Get("MatrixBot")
)

// The Bot struct holds parameters related to the bot
type Bot struct {
	receiveMessageChan chan events.ReceiveMessage
	sendMessageChan    chan events.SendMessage
	commandChan        chan events.Command
	server             string
	token              string
	pollingDone        chan bool

	pollingInterval time.Duration

	knownPlugins []plugins.Plugin

	// We are wasting a little bit of memory and keep maps in both directions
	knownRooms   map[string]string // mapping of Room to RoomID
	knownRoomIDs map[string]string // mapping of RoomID to Room

	nextBatch string // contains the next batch to fetch in sync
}

func (b Bot) apiCall(path string, method string, body string, auth bool) (r []byte, e error) {
	client := &http.Client{}

	req, err := http.NewRequest(method, b.server+"/_matrix"+path, strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	if auth == true {
		req.Header.Add("Authorization", "Bearer "+b.token)
	}
	req.Header.Add("Content-Type", "application/json")

	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(response.Body)
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

func (b *Bot) handlePolling() {
	b.callSync()
}

// CreateMatrixBot creates a new instance of a DiscordBot
func CreateMatrixBot(server string, username string, password string, token string) (*Bot, error) {
	log.Printf("MatrixBot is CREATING itself")
	b := Bot{server: server}
	if len(token) == 0 {
		token, err := b.login(username, password)
		if err != nil {
			return nil, err
		}
		b.token = token
	} else {
		// just use the provided access token
		b.token = token
	}

	b.connectToMatrixServer()

	b.pollingDone = make(chan bool)
	b.pollingInterval = 1000 * time.Millisecond

	b.receiveMessageChan = make(chan events.ReceiveMessage)
	b.sendMessageChan = make(chan events.SendMessage)
	b.commandChan = make(chan events.Command)

	b.knownRooms = make(map[string]string)
	b.knownRoomIDs = make(map[string]string)

	return &b, nil
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
