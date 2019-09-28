package fake

import (
	"github.com/torlenor/abylebotter/botinterface"
	"github.com/torlenor/abylebotter/events"
	"github.com/torlenor/abylebotter/logging"
	"github.com/torlenor/abylebotter/plugins"
)

var (
	log = logging.Get("FakeBot")
)

// The Bot struct holds parameters related to the bot
type Bot struct {
	receiveMessageChan chan events.ReceiveMessage
	sendMessageChan    chan events.SendMessage
	commandChan        chan events.Command

	pollingDone chan bool

	knownPlugins []plugins.Plugin
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

func (b *Bot) startBot(doneChannel chan struct{}) {
	defer close(doneChannel)
	// do some message polling or whatever until stopped

	for {
		select {
		case <-b.pollingDone:
			log.Println("polling stopped")
			return
		}
	}
}

// CreateFakeBot creates a new instance of a FakeBot
func CreateFakeBot() (*Bot, error) {
	log.Printf("FakeBot is CREATING itself")
	b := Bot{}

	b.pollingDone = make(chan bool)

	b.receiveMessageChan = make(chan events.ReceiveMessage)
	b.sendMessageChan = make(chan events.SendMessage)
	b.commandChan = make(chan events.Command)

	return &b, nil
}

func (b *Bot) startSendChannelReceiver() {
	for sendMsg := range b.sendMessageChan {
		switch sendMsg.Type {
		case events.MESSAGE:
			// do something
		case events.WHISPER:
			// do something
		default:
		}
	}
}

func (b *Bot) startCommandChannelReceiver() {
	for cmd := range b.commandChan {
		switch cmd.Command {
		case string("DemoCommand"):
			log.Println("Received DemoCommand with server name" + cmd.Payload)
		default:
			log.Println("Received unhandeled command" + cmd.Command)
		}
	}
}

// Start the Fake Bot
func (b *Bot) Start(doneChannel chan struct{}) {
	log.Println("FakeBot is STARTING")
	go b.startBot(doneChannel)
	go b.startSendChannelReceiver()
	go b.startCommandChannelReceiver()
}

// Stop the Fake Bot
func (b Bot) Stop() {
	log.Println("FakeBot is SHUTING DOWN")

	b.pollingDone <- true

	defer close(b.receiveMessageChan)
}

// Status returns the current status of FakeBot
func (b *Bot) Status() botinterface.BotStatus {
	return botinterface.BotStatus{Running: true}
}

// AddPlugin adds the give plugin to the current bot
func (b *Bot) AddPlugin(plugin plugins.Plugin) {
	plugin.ConnectChannels(b.GetReceiveMessageChannel(), b.GetSendMessageChannel(), b.GetCommandChannel())
	b.knownPlugins = append(b.knownPlugins, plugin)
}
