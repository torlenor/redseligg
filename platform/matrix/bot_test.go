package matrix

import (
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/torlenor/redseligg/commanddispatcher"
	"github.com/torlenor/redseligg/storage"
)

type mockAPI struct {
	letLoginFail   bool
	letAPICallFail bool
	apiResponse    string

	server    string
	authToken string

	loginCalled bool

	lastAPICallPath   string
	lastAPICallMethod string
	lastAPICallBody   string
	lastAPICallAuth   bool
}

func (api *mockAPI) call(path string, method string, body string, auth bool) (r []byte, e error) {
	api.lastAPICallPath = path
	api.lastAPICallMethod = method
	api.lastAPICallBody = body
	api.lastAPICallAuth = auth
	if api.letAPICallFail == true {
		return []byte(""), errors.New("Fake API call fail")
	}
	return []byte(api.apiResponse), nil
}

func (api *mockAPI) updateAuthToken(token string) {
	api.authToken = token
}

func (api *mockAPI) login(username string, password string) error {
	api.loginCalled = true
	if api.letLoginFail == true {
		return errors.New("Fake Login Fail")
	}
	return nil
}

func (api *mockAPI) reset() {
	api.letLoginFail = false
	api.letAPICallFail = false
	api.loginCalled = false
	api.authToken = ""
	api.lastAPICallPath = ""
	api.lastAPICallMethod = ""
	api.lastAPICallBody = ""
	api.lastAPICallAuth = false
}

func TestCreateMatrixBot(t *testing.T) {
	api := &mockAPI{server: "TEST_SERVER", authToken: "TEST_TOKEN"}
	api.reset()

	dispatcher := commanddispatcher.New("")
	storage := &storage.MockStorage{}

	// Login with user and password
	bot, err := createMatrixBotWithAPI(api, "TEST_USER", "TEST_PASS", dispatcher, storage)

	if bot == nil || err != nil {
		t.Fatalf("Could not create MatrixBot from username and password")
	}

	if api.loginCalled != true {
		t.Fatalf("Login not called even though user and password were provided")
	}

	// Fail login
	api.reset()
	api.letLoginFail = true

	bot, err = createMatrixBotWithAPI(api, "TEST_USER", "TEST_PASS", dispatcher, storage)

	if bot != nil || err == nil {
		t.Fatalf("Created Matrix bot even though login failed")
	}

	if api.loginCalled != true {
		t.Fatalf("Login not called even though user and password were provided")
	}
}

func TestMatrixBotPolling(t *testing.T) {
	api := &mockAPI{server: "TEST_SERVER", authToken: "TEST_TOKEN"}
	api.reset()

	dispatcher := commanddispatcher.New("")
	storage := &storage.MockStorage{}

	bot, err := createMatrixBotWithAPI(api, "TEST_USER", "TEST_PASS", dispatcher, storage)
	if bot == nil || err != nil {
		t.Fatalf("Could not create Matrix Bot")
	}

	// Polling without valid response from API has to fail
	err = bot.handlePolling()
	if api.lastAPICallPath != `/client/r0/sync?filter={"room":{"timeline":{"limit":1}}}` {
		t.Fatalf("handlePolling api call path wrong: %s", api.lastAPICallPath)
	}
	if api.lastAPICallMethod != `GET` {
		t.Fatalf("handlePolling api call method wrong: %s", api.lastAPICallMethod)
	}
	if api.lastAPICallBody != `{}` {
		t.Fatalf("handlePolling api call body wrong: %s", api.lastAPICallBody)
	}
	if api.lastAPICallAuth != true {
		t.Fatalf("handlePolling api call auth not set")
	}
	if err == nil {
		t.Fatalf("Initial sync somehow not failed even though it should")
	}

	// Polling with valid JSON response from API should not fail
	api.apiResponse = `{}`
	err = bot.handlePolling()
	if api.lastAPICallPath != `/client/r0/sync?filter={"room":{"timeline":{"limit":1}}}` {
		t.Fatalf("handlePolling api call path wrong: %s", api.lastAPICallPath)
	}
	if api.lastAPICallMethod != `GET` {
		t.Fatalf("handlePolling api call method wrong: %s", api.lastAPICallMethod)
	}
	if api.lastAPICallBody != `{}` {
		t.Fatalf("handlePolling api call body wrong: %s", api.lastAPICallBody)
	}
	if api.lastAPICallAuth != true {
		t.Fatalf("handlePolling api call auth not set")
	}
	if err != nil {
		t.Fatalf("Initial sync failed")
	}

}

func TestMatrixBotStartingAndStopping(t *testing.T) {
	api := &mockAPI{server: "TEST_SERVER", authToken: "TEST_TOKEN"}
	api.reset()

	dispatcher := commanddispatcher.New("")
	storage := &storage.MockStorage{}

	bot, err := createMatrixBotWithAPI(api, "TEST_USER", "TEST_PASS", dispatcher, storage)
	if bot == nil || err != nil {
		t.Fatalf("Could not create Matrix Bot")
	}

	go bot.Start()
	time.Sleep(time.Millisecond * 100)
	go bot.Stop()
	time.Sleep(time.Millisecond * 100)
}

func TestMatrixBotAddRemoveRoom(t *testing.T) {
	api := &mockAPI{server: "TEST_SERVER", authToken: "TEST_TOKEN"}
	api.reset()

	dispatcher := commanddispatcher.New("")
	storage := &storage.MockStorage{}

	bot, err := createMatrixBotWithAPI(api, "TEST_USER", "TEST_PASS", dispatcher, storage)
	if bot == nil || err != nil {
		t.Fatalf("Could not create Matrix Bot")
	}

	if len(bot.knownRooms) != 0 || len(bot.knownRoomIDs) != 0 {
		t.Fatalf("Initial room lists not empty")
	}

	bot.addKnownRoom("ID", "NAME")
	if len(bot.knownRooms) != 1 || len(bot.knownRoomIDs) != 1 {
		t.Fatalf("Room not successfully added, still zero rooms in list")
	}
	if bot.knownRooms["NAME"] != "ID" || bot.knownRoomIDs["ID"] != "NAME" {
		t.Fatalf("Room ID or Name not correctly added")
	}

	bot.removeKnownRoom("ID", "NAME")
	if len(bot.knownRooms) != 0 || len(bot.knownRoomIDs) != 0 {
		t.Fatalf("Room list not empty after removing the last room")
	}
}
