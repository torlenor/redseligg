package echoplugin

import (
	"testing"
	"time"

	"github.com/torlenor/abylebotter/botinterface"
	"github.com/torlenor/abylebotter/events"
)

func TestEcho(t *testing.T) {
	echoPlugin, err := CreateEchoPlugin()

	if err != nil {
		t.Fatalf("Could not create EchoPlugin")
	}

	var bot botinterface.MockBot
	bot.ReceiveMessageChan = make(chan events.ReceiveMessage)
	bot.SendMessageChan = make(chan events.SendMessage)

	bot.Reset()

	go bot.StartSendChannelReceiver()

	err = echoPlugin.ConnectChannels(bot.ReceiveMessageChan, bot.SendMessageChan)
	if err != nil {
		t.Fatalf("Error connecting channels")
	}

	msg := events.ReceiveMessage{Type: events.MESSAGE,
		ChannelID: "TEST_ChannelID_BEFORE_START",
		Content:   "TEST_MESSAGE_BEFORE_START"}
	bot.SendMessage(msg)

	time.Sleep(time.Millisecond * 20)

	if bot.LastSendMessage.ChannelID != "" ||
		bot.LastSendMessage.Content != "" ||
		bot.LastSendMessage.Type != events.UNKNOWN {
		t.Fatalf("mockBot already received a SendMessage request")
	}

	echoPlugin.Start()

	if echoPlugin.IsStarted() != true {
		t.Fatalf("EchoBot should report started now")
	}

	time.Sleep(time.Millisecond * 20)

	if bot.LastSendMessage.ChannelID != "" ||
		bot.LastSendMessage.Content != "" ||
		bot.LastSendMessage.Type != events.UNKNOWN {
		t.Fatalf("mockBot already received a SendMessage request")
	}

	msg = events.ReceiveMessage{Type: events.MESSAGE,
		ChannelID: "TEST_ChannelID_NO_ECHO",
		Content:   "TEST_MESSAGE_NO_ECHO"}
	bot.SendMessage(msg)

	time.Sleep(time.Millisecond * 20)

	if bot.LastSendMessage.ChannelID != "" ||
		bot.LastSendMessage.Content != "" ||
		bot.LastSendMessage.Type != events.UNKNOWN {
		t.Fatalf("mockBot received a SendMessage request even though message did not start with !echo")
	}

	msg = events.ReceiveMessage{Type: events.MESSAGE,
		ChannelID: "TEST_ChannelID",
		Content:   "!echo TEST_MESSAGE"}
	bot.SendMessage(msg)

	time.Sleep(time.Millisecond * 20)

	if bot.LastSendMessage.ChannelID != "TEST_ChannelID" ||
		bot.LastSendMessage.Content != "TEST_MESSAGE" ||
		bot.LastSendMessage.Type != events.MESSAGE {
		t.Fatalf("mockBot did not receive a SendMessage request even though the EchoPlugin should have echoed it")
	}

	bot.Reset()
	go bot.StartSendChannelReceiver()

	msg = events.ReceiveMessage{Type: events.WHISPER,
		ChannelID: "TEST_ChannelID_WHISPER",
		Content:   "!echo TEST_WHISPER"}
	bot.SendMessage(msg)

	time.Sleep(time.Millisecond * 20)

	if bot.LastSendMessage.ChannelID != "TEST_ChannelID_WHISPER" ||
		bot.LastSendMessage.Content != "TEST_WHISPER" ||
		bot.LastSendMessage.Type != events.WHISPER {
		t.Fatalf("EchoBot did not echo WHISPER")
	}

	bot.Reset()
	go bot.StartSendChannelReceiver()
	echoPlugin.SetOnlyOnWhisper(true)

	msg = events.ReceiveMessage{Type: events.MESSAGE,
		ChannelID: "TEST_ChannelID",
		Content:   "!echo TEST_YET_ANOTHER_MESSAGE"}
	bot.SendMessage(msg)

	time.Sleep(time.Millisecond * 20)

	if bot.LastSendMessage.ChannelID != "" ||
		bot.LastSendMessage.Content != "" ||
		bot.LastSendMessage.Type != events.UNKNOWN {
		t.Fatalf("EchoBot echoed a MESSAGE even though it is set to Whisper Only")
	}

	msg = events.ReceiveMessage{Type: events.WHISPER,
		ChannelID: "TEST_ChannelID_ANOTHER_WHISPER",
		Content:   "!echo TEST_ANOTHER_WHISPER"}
	bot.SendMessage(msg)

	time.Sleep(time.Millisecond * 20)

	if bot.LastSendMessage.ChannelID != "TEST_ChannelID_ANOTHER_WHISPER" ||
		bot.LastSendMessage.Content != "TEST_ANOTHER_WHISPER" ||
		bot.LastSendMessage.Type != events.WHISPER {
		t.Fatalf("EchoBot did not echo WHISPER")
	}

	bot.Reset()
	echoPlugin.Stop()

	if echoPlugin.IsStarted() != false {
		t.Fatalf("EchoBot should report stopped now")
	}

	msg = events.ReceiveMessage{Type: events.MESSAGE,
		ChannelID: "TEST_ChannelID",
		Content:   "!echo TEST_YET_YET_ANOTHER_MESSAGE"}
	bot.SendMessage(msg)

	msg = events.ReceiveMessage{Type: events.WHISPER,
		ChannelID: "TEST_ChannelID_WHISPER",
		Content:   "!echo TEST_YET_ANOTHER_WHISPER"}
	bot.SendMessage(msg)

	time.Sleep(time.Millisecond * 20)

	if bot.LastSendMessage.ChannelID != "" ||
		bot.LastSendMessage.Content != "" ||
		bot.LastSendMessage.Type != events.UNKNOWN {
		t.Fatalf("EchoBot echoed something even though it is stopped")
	}

	bot.Reset()
	echoPlugin.Start()
	close(bot.ReceiveMessageChan)

	time.Sleep(time.Millisecond * 20)

	if echoPlugin.IsStarted() != false {
		t.Fatalf("EchoBot should have been stopped automatically on receiveChannel close")
	}
}

func TestEchoPlugin_GetName(t *testing.T) {
	want := "EchoPlugin"
	p := &EchoPlugin{}
	if got := p.GetName(); got != want {
		t.Errorf("EchoPlugin.GetName() = %v, want %v", got, want)
	}
}
