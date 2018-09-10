package matrix

import (
	"testing"

	"github.com/pkg/errors"
)

type mockAPI struct {
	letLoginFail bool

	server    string
	authToken string

	loginCalled                 bool
	connectToMatrixServerCalled bool
}

func (api *mockAPI) call(path string, method string, body string, auth bool) (r []byte, e error) {
	return []byte(""), nil
}

func (api *mockAPI) updateAuthToken(token string) {
	api.authToken = token
}

func (api *mockAPI) connectToMatrixServer() error {
	api.connectToMatrixServerCalled = true
	return nil
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
	api.loginCalled = false
	api.connectToMatrixServerCalled = false
	api.authToken = ""
}

func TestCreateMatrixBot(t *testing.T) {
	api := &mockAPI{server: "TEST_SERVER", authToken: "TEST_TOKEN"}
	api.reset()

	// Login with user and password
	bot, err := createMatrixBotWithAPI(api, "TEST_USER", "TEST_PASS", "")

	if bot == nil || err != nil {
		t.Fatalf("Could not create MatrixBot from username and password")
	}

	if api.loginCalled != true {
		t.Fatalf("Login not called even though user and password were provided")
	}

	if api.connectToMatrixServerCalled != true {
		t.Fatalf("Not connected to Matrix Server")
	}

	// Fail login
	api.reset()
	api.letLoginFail = true

	bot, err = createMatrixBotWithAPI(api, "TEST_USER", "TEST_PASS", "")

	if bot != nil || err == nil {
		t.Fatalf("Created Matrix bot even though login failed")
	}

	if api.loginCalled != true {
		t.Fatalf("Login not called even though user and password were provided")
	}

	if api.connectToMatrixServerCalled != false {
		t.Fatalf("Tried to connect to Matrix server even though login failed")
	}

	// No login, but token provided
	api.reset()

	bot, err = createMatrixBotWithAPI(api, "", "", "TEST_TOKEN")

	if bot == nil || err != nil {
		t.Fatalf("Could not create MatrixBot from token")
	}

	if api.loginCalled != false {
		t.Fatalf("Login called even though token was provided")
	}

	if api.connectToMatrixServerCalled != true {
		t.Fatalf("Not connected to Matrix Server")
	}

	if api.authToken != "TEST_TOKEN" {
		t.Fatalf("Auth token was not updated")
	}

}
