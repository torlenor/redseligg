package discord

import (
	"fmt"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/torlenor/abylebotter/ws"
)

func Test_SendIdent(t *testing.T) {
	ws := &ws.MockClient{}
	defer ws.Stop()

	testToken := "TESTTOKEN.12345"

	expectedIdent := []byte(`{"op": 2,
			"d": {
				"token": "` + testToken + `",
				"properties": {},
				"compress": false,
				"large_threshold": 250
			}
}`)

	err := sendIdent(testToken, ws)
	if err != nil {
		t.Fatalf("Sending Ident failed")
	}

	if string(ws.LastSendMessageData) != string(expectedIdent) {
		t.Fatalf("Bad message was sent")
	}
	if ws.LastSendMessageType != websocket.TextMessage {
		t.Fatalf("Bad message type was used")
	}
}

func Test_SendIdentFail(t *testing.T) {
	ws := &ws.MockClient{}
	defer ws.Stop()

	testToken := "TESTTOKEN.12345"

	expectedIdent := []byte(`{"op": 2,
			"d": {
				"token": "` + testToken + `",
				"properties": {},
				"compress": false,
				"large_threshold": 250
			}
}`)

	ws.ReturnError = fmt.Errorf("Some error")

	err := sendIdent(testToken, ws)
	if err == nil {
		t.Fatalf("Sending Ident did not fail")
	}

	if string(ws.LastSendMessageData) != string(expectedIdent) {
		t.Fatalf("Bad message was sent")
	}
	if ws.LastSendMessageType != websocket.TextMessage {
		t.Fatalf("Bad message type was used")
	}
}
