package discord

import (
	"fmt"
	"testing"

	"git.abyle.org/redseligg/botorchestrator/botconfig"

	"github.com/torlenor/abylebotter/storage"
	"github.com/torlenor/abylebotter/webclient"
	"github.com/torlenor/abylebotter/ws"
)

func Test_CreateDiscordBot(t *testing.T) {
	ws := &ws.MockClient{}
	storage := storage.MockStorage{}
	api := webclient.NewMock()

	expectedAPICallPath := "/gateway"
	expectedAPICallMethod := "GET"
	expectedAPICallBody := ""
	expectedWebSocketGatewayURL := "ws://something"

	api.ReturnOnCall = webclient.APIResponse{
		Body: []byte(`{"url": "` + expectedWebSocketGatewayURL + `"}`),
	}

	cfg := botconfig.DiscordConfig{}

	bot, err := CreateDiscordBotWithAPI(api, &storage, cfg, ws)
	if err != nil {
		t.Fatalf("Creating the bot should not have failed")
	}
	if api.LastCallPath != expectedAPICallPath {
		t.Fatalf("Did not trigger correct API Call, path wanted %s, got %s", expectedAPICallPath, api.LastCallPath)
	}
	if api.LastCallMethod != expectedAPICallMethod {
		t.Fatalf("Did not trigger correct API Call, method wanted %s, got %s", expectedAPICallMethod, api.LastCallMethod)
	}
	if api.LastCallBody != expectedAPICallBody {
		t.Fatalf("Did not trigger correct API Call, body wanted %s, got %s", expectedAPICallBody, api.LastCallBody)
	}
	if bot.gatewayURL != expectedWebSocketGatewayURL {
		t.Fatalf("Websocket Gateway URL wrong, wanted: %s, got: %s", expectedWebSocketGatewayURL, bot.gatewayURL)
	}
	if ws.WasDialCalled != false {
		t.Fatalf("Should not have dialed the WebSocket Gateway")
	}

	api.Reset()
	api.ReturnOnCallError = fmt.Errorf("Some error")
	bot, err = CreateDiscordBotWithAPI(api, &storage, cfg, ws)
	if err == nil {
		t.Fatalf("Creating the bot should have failed")
	}
	if api.LastCallPath != expectedAPICallPath {
		t.Fatalf("Did not trigger correct API Call, path wanted %s, got %s", expectedAPICallPath, api.LastCallPath)
	}
	if api.LastCallMethod != expectedAPICallMethod {
		t.Fatalf("Did not trigger correct API Call, method wanted %s, got %s", expectedAPICallMethod, api.LastCallMethod)
	}
	if api.LastCallBody != expectedAPICallBody {
		t.Fatalf("Did not trigger correct API Call, body wanted %s, got %s", expectedAPICallBody, api.LastCallBody)
	}
}
