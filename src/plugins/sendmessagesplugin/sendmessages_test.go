package sendmessagesplugin

import (
	"events"
	"net/http"
	"testing"
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

// func TestSendMessagesPlugin_handleSendMessage(t *testing.T) {
// 	type fields struct {
// 		log               *logrus.Entry
// 		botSendChannel    chan events.SendMessage
// 		botCommandChannel chan events.Command
// 		isStarted         bool
// 	}
// 	type args struct {
// 		ident   string
// 		content string
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := &SendMessagesPlugin{
// 				log:               tt.fields.log,
// 				botSendChannel:    tt.fields.botSendChannel,
// 				botCommandChannel: tt.fields.botCommandChannel,
// 				isStarted:         tt.fields.isStarted,
// 			}
// 			p.handleSendMessage(tt.args.ident, tt.args.content)
// 		})
// 	}
// }

// func TestSendMessagesPlugin_sendMessage(t *testing.T) {
// 	type fields struct {
// 		log               *logrus.Entry
// 		botSendChannel    chan events.SendMessage
// 		botCommandChannel chan events.Command
// 		isStarted         bool
// 	}
// 	type args struct {
// 		w http.ResponseWriter
// 		r *http.Request
// 	}
// 	tests := []struct {
// 		name   string
// 		fields fields
// 		args   args
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			p := &SendMessagesPlugin{
// 				log:               tt.fields.log,
// 				botSendChannel:    tt.fields.botSendChannel,
// 				botCommandChannel: tt.fields.botCommandChannel,
// 				isStarted:         tt.fields.isStarted,
// 			}
// 			p.sendMessage(tt.args.w, tt.args.r)
// 		})
// 	}
// }

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
	commandCHannel := make(chan events.Command)

	p := &SendMessagesPlugin{}
	if err := p.ConnectChannels(receiveChannel, sendChannel, commandCHannel); err != nil {
		t.Errorf("SendMessagesPlugin.ConnectChannels() error = %v, wantErr %v", err, false)
	}
	if p.botSendChannel != sendChannel {
		t.Errorf("SendMessagesPlugin.ConnectChannels() sendChannel not connected")
	}
	if p.botCommandChannel != commandCHannel {
		t.Errorf("SendMessagesPlugin.ConnectChannels() commandCHannel not connected")
	}
}
