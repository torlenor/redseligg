package plugins

import (
	"events"
	"testing"
	"time"
)

type mockBot struct {
	receiveMessageChan chan events.ReceiveMessage
	sendMessageChan    chan events.SendMessage
	commandChan        chan events.Command

	lastSendMessage     events.SendMessage
	lastReceivedCommand events.Command
}

func (b *mockBot) startSendChannelReceiver(done chan bool) {
	b.lastSendMessage = <-b.sendMessageChan
	done <- true
}

func (b *mockBot) startCommandChannelReceiver(done chan bool) {
	b.lastReceivedCommand = <-b.commandChan
	done <- true
}

func (b *mockBot) reset() {
	b.lastSendMessage = events.SendMessage{Type: events.UNKNOWN, Ident: "", Content: ""}
	b.lastReceivedCommand = events.Command{Command: "", Payload: ""}
}

func (b *mockBot) sendMessage(msg events.ReceiveMessage) {
	select {
	case b.receiveMessageChan <- msg:
	default:
	}
}

func TestEcho(t *testing.T) {
	echoPlugin, err := CreateEchoPlugin()

	if err != nil {
		t.Fatalf("Could not create EchoPlugin")
	}

	var bot mockBot
	bot.receiveMessageChan = make(chan events.ReceiveMessage)
	bot.sendMessageChan = make(chan events.SendMessage)
	bot.commandChan = make(chan events.Command)

	doneSend := make(chan bool)
	defer close(doneSend)

	bot.reset()

	go bot.startSendChannelReceiver(doneSend)

	err = echoPlugin.ConnectChannels(bot.receiveMessageChan, bot.sendMessageChan, bot.commandChan)
	if err != nil {
		t.Fatalf("Error connecting channels")
	}

	msg := events.ReceiveMessage{Type: events.MESSAGE,
		Ident:   "TEST_IDENT_BEFORE_START",
		Content: "TEST_MESSAGE_BEFORE_START"}
	bot.sendMessage(msg)

	time.Sleep(time.Millisecond * 20)

	if bot.lastSendMessage.Ident != "" ||
		bot.lastSendMessage.Content != "" ||
		bot.lastSendMessage.Type != events.UNKNOWN {
		t.Fatalf("mockBot already received a SendMessage request")
	}

	echoPlugin.Start()

	if echoPlugin.IsStarted() != true {
		t.Fatalf("EchoBot should report started now")
	}

	time.Sleep(time.Millisecond * 20)

	if bot.lastSendMessage.Ident != "" ||
		bot.lastSendMessage.Content != "" ||
		bot.lastSendMessage.Type != events.UNKNOWN {
		t.Fatalf("mockBot already received a SendMessage request")
	}

	msg = events.ReceiveMessage{Type: events.MESSAGE,
		Ident:   "TEST_IDENT_NO_ECHO",
		Content: "TEST_MESSAGE_NO_ECHO"}
	bot.sendMessage(msg)

	time.Sleep(time.Millisecond * 20)

	if bot.lastSendMessage.Ident != "" ||
		bot.lastSendMessage.Content != "" ||
		bot.lastSendMessage.Type != events.UNKNOWN {
		t.Fatalf("mockBot received a SendMessage request even though message did not start with !echo")
	}

	msg = events.ReceiveMessage{Type: events.MESSAGE,
		Ident:   "TEST_IDENT",
		Content: "!echo TEST_MESSAGE"}
	bot.sendMessage(msg)

	<-doneSend
	if bot.lastSendMessage.Ident != "TEST_IDENT" ||
		bot.lastSendMessage.Content != "TEST_MESSAGE" ||
		bot.lastSendMessage.Type != events.MESSAGE {
		t.Fatalf("mockBot did not receive a SendMessage request even though the EchoPlugin should have echoed it")
	}

	bot.reset()
	go bot.startSendChannelReceiver(doneSend)

	msg = events.ReceiveMessage{Type: events.WHISPER,
		Ident:   "TEST_IDENT_WHISPER",
		Content: "!echo TEST_WHISPER"}
	bot.sendMessage(msg)

	<-doneSend
	if bot.lastSendMessage.Ident != "TEST_IDENT_WHISPER" ||
		bot.lastSendMessage.Content != "TEST_WHISPER" ||
		bot.lastSendMessage.Type != events.WHISPER {
		t.Fatalf("EchoBot did not echo WHISPER")
	}

	bot.reset()
	go bot.startSendChannelReceiver(doneSend)
	echoPlugin.SetOnlyOnWhisper(true)

	msg = events.ReceiveMessage{Type: events.MESSAGE,
		Ident:   "TEST_IDENT",
		Content: "!echo TEST_YET_ANOTHER_MESSAGE"}
	bot.sendMessage(msg)

	time.Sleep(time.Millisecond * 20)

	if bot.lastSendMessage.Ident != "" ||
		bot.lastSendMessage.Content != "" ||
		bot.lastSendMessage.Type != events.UNKNOWN {
		t.Fatalf("EchoBot echoed a MESSAGE even though it is set to Whisper Only")
	}

	msg = events.ReceiveMessage{Type: events.WHISPER,
		Ident:   "TEST_IDENT_ANOTHER_WHISPER",
		Content: "!echo TEST_ANOTHER_WHISPER"}
	bot.sendMessage(msg)

	<-doneSend
	if bot.lastSendMessage.Ident != "TEST_IDENT_ANOTHER_WHISPER" ||
		bot.lastSendMessage.Content != "TEST_ANOTHER_WHISPER" ||
		bot.lastSendMessage.Type != events.WHISPER {
		t.Fatalf("EchoBot did not echo WHISPER")
	}

	bot.reset()
	echoPlugin.Stop()

	if echoPlugin.IsStarted() != false {
		t.Fatalf("EchoBot should report stopped now")
	}

	msg = events.ReceiveMessage{Type: events.MESSAGE,
		Ident:   "TEST_IDENT",
		Content: "!echo TEST_YET_YET_ANOTHER_MESSAGE"}
	bot.sendMessage(msg)

	msg = events.ReceiveMessage{Type: events.WHISPER,
		Ident:   "TEST_IDENT_WHISPER",
		Content: "!echo TEST_YET_ANOTHER_WHISPER"}
	bot.sendMessage(msg)

	time.Sleep(time.Millisecond * 20)

	if bot.lastSendMessage.Ident != "" ||
		bot.lastSendMessage.Content != "" ||
		bot.lastSendMessage.Type != events.UNKNOWN {
		t.Fatalf("EchoBot echoed something even though it is stopped")
	}

	bot.reset()
	echoPlugin.Start()
	close(bot.receiveMessageChan)

	time.Sleep(time.Millisecond * 20)

	if echoPlugin.IsStarted() != false {
		t.Fatalf("EchoBot should have been stopped automatically on receiveChannel close")
	}

}
