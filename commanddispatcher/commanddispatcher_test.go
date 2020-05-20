package commanddispatcher

import (
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

	dispatcher.Unregister(expectedCommand)
	assert.Equal(0, len(dispatcher.receivers))
}
