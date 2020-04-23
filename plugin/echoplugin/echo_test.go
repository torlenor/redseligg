package echoplugin

import (
	"testing"

	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/plugin"

	"github.com/stretchr/testify/assert"
)

func TestEchoPlugin_OnPost(t *testing.T) {
	assert := assert.New(t)

	p := EchoPlugin{}
	assert.Equal(nil, p.API)
	assert.Equal(false, p.onlyOnWhisper)

	api := plugin.MockAPI{}
	storage := plugin.MockStorageAPI{}
	p.SetAPI(&api, &storage)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.Content = "!echo"
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.Content = "!echo MESSAGE CONTENT"
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	p.SetOnlyOnWhisper(true)
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.IsPrivate = true
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
}
