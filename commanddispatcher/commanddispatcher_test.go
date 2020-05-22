package commanddispatcher

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/torlenor/redseligg/model"
)

type mockCommandReceiver struct {
	lastReceivedCmd     string
	lastReceivedContent string
	lastReceivedPost    model.Post
}

func (m *mockCommandReceiver) OnCommand(cmd string, content string, post model.Post) {
	m.lastReceivedCmd = cmd
	m.lastReceivedContent = content
	m.lastReceivedPost = post
}

func TestCommandDispatcher(t *testing.T) {
	assert := assert.New(t)

	dispatcher := New("")
	assert.Equal(defaultCallPrefix, dispatcher.callPrefix)

	expectedCallPrefix := "~"
	dispatcher = New(expectedCallPrefix)
	assert.Equal(expectedCallPrefix, dispatcher.callPrefix)

	expectedCommand := "someCommand"
	receiver := &mockCommandReceiver{}
	assert.Equal(0, len(dispatcher.receivers))
	dispatcher.Register(expectedCommand, receiver)
	assert.Equal(1, len(dispatcher.receivers))
	assert.Equal(receiver, dispatcher.receivers[expectedCommand])

	assert.Equal("", receiver.lastReceivedCmd)
	assert.Equal(model.Post{}, receiver.lastReceivedPost)

	expectedPost := model.Post{
		ChannelID: "some id",
		Channel:   "some channel",

		Content: "!otherCommand",
	}
	dispatcher.OnPost(expectedPost)
	assert.Equal("", receiver.lastReceivedCmd)
	assert.Equal("", receiver.lastReceivedContent)
	assert.Equal(model.Post{}, receiver.lastReceivedPost)

	expectedPost.Content = "!" + expectedCommand
	dispatcher.OnPost(expectedPost)
	assert.Equal("", receiver.lastReceivedCmd)
	assert.Equal("", receiver.lastReceivedContent)
	assert.Equal(model.Post{}, receiver.lastReceivedPost)

	expectedPost.Content = "~" + expectedCommand
	dispatcher.OnPost(expectedPost)
	assert.Equal(expectedCommand, receiver.lastReceivedCmd)
	assert.Equal("", receiver.lastReceivedContent)
	assert.Equal(expectedPost, receiver.lastReceivedPost)

	expectedContent := "some content"
	expectedPost.Content = "~" + expectedCommand + " " + expectedContent
	dispatcher.OnPost(expectedPost)
	assert.Equal(expectedCommand, receiver.lastReceivedCmd)
	assert.Equal(expectedContent, receiver.lastReceivedContent)
	assert.Equal(expectedPost, receiver.lastReceivedPost)

	assert.Equal(fmt.Sprintf("The following commands are available: ~%s\nNote: Some of them are only available for mods.", expectedCommand), dispatcher.getHelpText())

	dispatcher.Unregister(expectedCommand)
	assert.Equal(0, len(dispatcher.receivers))
}

func TestCommandDispatcher_GetCallPrefix(t *testing.T) {
	type fields struct {
		callPrefix string
		receivers  map[string]receiver
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name:   "prefix !",
			fields: fields{callPrefix: "!"},
			want:   "!",
		},
		{
			name:   "prefix abc",
			fields: fields{callPrefix: "abc"},
			want:   "abc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandDispatcher{
				callPrefix: tt.fields.callPrefix,
				receivers:  tt.fields.receivers,
			}
			if got := c.GetCallPrefix(); got != tt.want {
				t.Errorf("CommandDispatcher.GetCallPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommandDispatcher_IsHelp(t *testing.T) {
	type fields struct {
		callPrefix string
		receivers  map[string]receiver
	}
	type args struct {
		post model.Post
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
		want1  string
	}{
		{
			name:   "Default call prefix: is help",
			fields: fields{callPrefix: "!"},
			args: args{
				post: model.Post{Content: "!help"},
			},
			want:  true,
			want1: "The following commands are available: \nNote: Some of them are only available for mods.",
		},
		{
			name:   "Default call prefix: is not help",
			fields: fields{callPrefix: "!"},
			args: args{
				post: model.Post{Content: "!nothelp"},
			},
			want:  false,
			want1: "",
		},
		{
			name:   "Default call prefix: also not help",
			fields: fields{callPrefix: "!"},
			args: args{
				post: model.Post{Content: "!helpabcabc"},
			},
			want:  false,
			want1: "",
		},
		{
			name:   "Other call prefix: is help",
			fields: fields{callPrefix: "abc"},
			args: args{
				post: model.Post{Content: "abchelp"},
			},
			want:  true,
			want1: "The following commands are available: \nNote: Some of them are only available for mods.",
		},
		{
			name:   "Other call prefix: is not help",
			fields: fields{callPrefix: "abc"},
			args: args{
				post: model.Post{Content: "abcnothelp"},
			},
			want:  false,
			want1: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CommandDispatcher{
				callPrefix: tt.fields.callPrefix,
				receivers:  tt.fields.receivers,
			}
			got, got1 := c.IsHelp(tt.args.post)
			if got != tt.want {
				t.Errorf("CommandDispatcher.IsHelp() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CommandDispatcher.IsHelp() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
