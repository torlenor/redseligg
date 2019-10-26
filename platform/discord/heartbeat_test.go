package discord

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestDiscordHeartBeatSender(t *testing.T) {
	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	testSeqNumber := 1
	expectedHeartbeat := []byte(`{"op":1,"d":` + strconv.Itoa(testSeqNumber) + `}`)

	var dhbs discordHeartBeatSender
	dhbs.ws = ws
	dhbs.sendHeartBeat(testSeqNumber)

	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if string(p) != string(expectedHeartbeat) {
		t.Fatalf("bad message")
	}
}

type mockHeartBeatSender struct {
	lastHeartBeatSent string
}

func (hbs *mockHeartBeatSender) sendHeartBeat(seqNumber int) error {
	hbs.lastHeartBeatSent = `{"op":1,"d":` + strconv.Itoa(seqNumber) + `}`
	return nil
}

func TestHeartBeat(t *testing.T) {
	var mockSender = &mockHeartBeatSender{}

	stopHeartBeat := make(chan struct{})
	seqNumberChan := make(chan int)

	// If anybody knows a better way in golang than using sleeps, please tell me

	go heartBeat(10, mockSender, stopHeartBeat, seqNumberChan)

	for i := 1; i <= 10; i++ {
		seqNumberChan <- i
		time.Sleep(time.Millisecond * 20)
		if mockSender.lastHeartBeatSent != string(`{"op":1,"d":`+strconv.Itoa(i)+`}`) {
			t.Fatalf("bad heartbeat message received for seqNumber %d", i)
		}
	}

	close(stopHeartBeat)

	time.Sleep(time.Millisecond * 20)
}
