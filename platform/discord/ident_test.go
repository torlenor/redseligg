package discord

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/torlenor/abylebotter/ws"
)

var upgrader = websocket.Upgrader{}

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			break
		}
		err = c.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func TestSendIdent(t *testing.T) {
	// TODO: Refactor that so we are using a mock and not real websocket connections
	// Create test server with the echo handler.
	s := httptest.NewServer(http.HandlerFunc(echo))
	defer s.Close()

	u := "ws" + strings.TrimPrefix(s.URL, "http")

	// Connect to the server
	ws := ws.NewClient()
	defer ws.Stop()
	ws.Dial(u)

	testToken := "TESTTOKEN.12345"

	expectedIdent := []byte(`{"op": 2,
			"d": {
				"token": "` + testToken + `",
				"properties": {},
				"compress": false,
				"large_threshold": 250
			}
}`)

	// Send message to server, read response and check to see if it's what we expect.
	sendIdent(testToken, ws)

	_, p, err := ws.ReadMessage()
	if err != nil {
		t.Fatalf("%v", err)
	}

	if string(p) != string(expectedIdent) {
		t.Fatalf("bad message")
	}
}
