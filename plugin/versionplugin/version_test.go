package versionplugin

import (
	"testing"

	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/plugin"

	"github.com/stretchr/testify/assert"
)

func TestEchoPlugin_OnPost(t *testing.T) {
	assert := assert.New(t)

	p := VersionPlugin{}
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	p.SetAPI(&api)

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
	postToPlugin.Content = "!version"
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   api.GetVersion(),
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}
