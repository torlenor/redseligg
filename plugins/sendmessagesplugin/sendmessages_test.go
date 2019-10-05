package sendmessagesplugin

import (
	"net/http"
	"testing"

	"github.com/torlenor/abylebotter/events"
)

type mockAPI struct {
	path     string
	function func(http.ResponseWriter, *http.Request)
}

func (ma *mockAPI) AttachModulePost(path string, f func(http.ResponseWriter, *http.Request)) {
	ma.path = path
	ma.function = f
}

func TestCreateSendMessagesPlugin(t *testing.T) {
	got, err := CreateSendMessagesPlugin()
	if err != nil {
		t.Errorf("CreateSendMessagesPlugin() error = %v", err)
		return
	}
	want := false
	if got.isStarted != want {
		t.Errorf("CreateSendMessagesPlugin() isStarted = %v, want = %v", got.isStarted, want)
		return
	}
}

func TestSendMessagesPlugin_RegisterToRestAPI(t *testing.T) {
	mapi := &mockAPI{}
	registerToAPI = mapi.AttachModulePost

	p := &SendMessagesPlugin{}
	p.RegisterToRestAPI()

	want := "/plugins/sendmessages"
	if got := mapi.path; got != want {
		t.Errorf("SendMessagesPlugin.RegisterToRestAPI() path = %v, want %v", got, want)
	}
}

func TestSendMessagesPlugin_Start_Stop(t *testing.T) {

	p := &SendMessagesPlugin{}
	want := false
	if got := p.IsStarted(); got != want {
		t.Errorf("Before SendMessagesPlugin.Start() call IsStarted() = %v, want %v", got, want)
	}
	p.Start()

	want = true
	if got := p.IsStarted(); got != want {
		t.Errorf("After SendMessagesPlugin.Start() call IsStarted() = %v, want %v", got, want)
	}

	p.Stop()
	want = false
	if got := p.IsStarted(); got != want {
		t.Errorf("After SendMessagesPlugin.Stop() call IsStarted() = %v, want %v", got, want)
	}
}

func TestSendMessagesPlugin_ConnectChannels(t *testing.T) {
	receiveChannel := make(chan events.ReceiveMessage)
	sendChannel := make(chan events.SendMessage)
	p := &SendMessagesPlugin{}
	if err := p.ConnectChannels(receiveChannel, sendChannel); err != nil {
		t.Errorf("SendMessagesPlugin.ConnectChannels() error = %v, wantErr %v", err, false)
	}
	if p.botSendChannel != sendChannel {
		t.Errorf("SendMessagesPlugin.ConnectChannels() sendChannel not connected")
	}
}

func TestSendMessagesPlugin_GetName(t *testing.T) {
	want := "SendMessagesPlugin"
	p := &SendMessagesPlugin{}
	if got := p.GetName(); got != want {
		t.Errorf("SendMessagesPlugin.GetName() = %v, want %v", got, want)
	}
}
