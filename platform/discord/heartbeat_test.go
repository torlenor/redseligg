package discord

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/torlenor/abylebotter/ws"
)

type onFailMock struct {
	OnFailCalled bool
}

func (f *onFailMock) onFail() {
	f.OnFailCalled = true
}

func TestDiscordHeartBeatSender(t *testing.T) {
	// Connect to the server
	ws := &ws.MockClient{}
	defer ws.Stop()

	testSeqNumber := 1
	expectedHeartbeat := []byte(`{"op":1,"d":` + strconv.Itoa(testSeqNumber) + `}`)

	var dhbs discordHeartBeatSender
	dhbs.ws = ws
	dhbs.sendHeartBeat(testSeqNumber)

	if string(ws.LastSendMessageData) != string(expectedHeartbeat) {
		t.Fatalf("bad message")
	}
}

type mockHeartBeatSender struct {
	lastHeartBeatSent string

	ReturnError error
}

func (hbs *mockHeartBeatSender) sendHeartBeat(seqNumber int) error {
	hbs.lastHeartBeatSent = `{"op":1,"d":` + strconv.Itoa(seqNumber) + `}`
	return hbs.ReturnError
}

func Test_HeartBeat(t *testing.T) {
	onFailHandler := onFailMock{}
	mockSender := &mockHeartBeatSender{}

	stopHeartBeat := make(chan bool)
	seqNumberChan := make(chan int)

	// If anybody knows a better way in golang than using sleeps, please tell me

	go heartBeat(10, mockSender, stopHeartBeat, seqNumberChan, onFailHandler.onFail)

	for i := 1; i <= 10; i++ {
		seqNumberChan <- i
		time.Sleep(time.Millisecond * 20)
		if mockSender.lastHeartBeatSent != string(`{"op":1,"d":`+strconv.Itoa(i)+`}`) {
			t.Fatalf("bad heartbeat message received for seqNumber %d", i)
		}
	}

	if onFailHandler.OnFailCalled {
		t.Fatalf("onFail was called")
	}

	close(stopHeartBeat)

	time.Sleep(time.Millisecond * 20)
}

func Test_HeartBeatSendFail(t *testing.T) {
	onFailHandler := onFailMock{}
	mockSender := &mockHeartBeatSender{}

	stopHeartBeat := make(chan bool)
	seqNumberChan := make(chan int)

	mockSender.ReturnError = fmt.Errorf("Some error")

	// If anybody knows a better way in golang than using sleeps, please tell me

	go heartBeat(10, mockSender, stopHeartBeat, seqNumberChan, onFailHandler.onFail)

	i := 1
	seqNumberChan <- i
	time.Sleep(time.Millisecond * 20)
	if mockSender.lastHeartBeatSent != string(`{"op":1,"d":`+strconv.Itoa(i)+`}`) {
		t.Fatalf("bad heartbeat message received for seqNumber %d", i)
	}

	if !onFailHandler.OnFailCalled {
		t.Fatalf("onFail was not called")
	}

	close(stopHeartBeat)

	time.Sleep(time.Millisecond * 20)
}
