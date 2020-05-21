package matrix

import (
	"time"

	"github.com/torlenor/redseligg/botconfig"
	"github.com/torlenor/redseligg/commanddispatcher"

	"github.com/torlenor/redseligg/logging"
	"github.com/torlenor/redseligg/platform"
	"github.com/torlenor/redseligg/plugin"
)

var (
	log = logging.Get("MatrixBot")
)

// The Bot struct holds parameters related to the bot
type Bot struct {
	platform.BotImpl
	api api

	dispatcher *commanddispatcher.CommandDispatcher

	pollingDone chan bool

	pollingInterval time.Duration

	plugins []plugin.Hooks

	// We are wasting a little bit of memory and keep maps in both directions
	knownRooms   map[string]string // mapping of Room to RoomID
	knownRoomIDs map[string]string // mapping of RoomID to Room

	nextBatch string // contains the next batch to fetch in sync
}

// The createMatrixBotWithAPI creates a new instance of a MatrixBot using the api interface api
func createMatrixBotWithAPI(api api, username string, password string, commandDispatcher *commanddispatcher.CommandDispatcher) (*Bot, error) {
	log.Printf("MatrixBot is CREATING itself")
	b := Bot{
		api:        api,
		dispatcher: commandDispatcher,
	}

	err := b.api.login(username, password)
	if err != nil {
		return nil, err
	}

	b.pollingDone = make(chan bool)
	b.pollingInterval = 1000 * time.Millisecond

	b.knownRooms = make(map[string]string)
	b.knownRoomIDs = make(map[string]string)

	return &b, nil
}

// CreateMatrixBot creates a new instance of a DiscordBot
func CreateMatrixBot(cfg botconfig.MatrixConfig, commandDispatcher *commanddispatcher.CommandDispatcher) (*Bot, error) {
	api := &matrixAPI{server: cfg.Server}
	return createMatrixBotWithAPI(api, cfg.Username, cfg.Password, commandDispatcher)
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

// GetInfo returns information about the Bot
func (b *Bot) GetInfo() platform.BotInfo {
	return platform.BotInfo{
		BotID:    "",
		Platform: "Matrix",
		Healthy:  true,
		Plugins:  []platform.PluginInfo{},
	}
}
