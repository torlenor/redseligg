package versionplugin

import (
	"testing"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/plugin"

	"github.com/stretchr/testify/assert"
)

func TestEchoPlugin_OnCommand(t *testing.T) {
	assert := assert.New(t)

	p := VersionPlugin{}
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	p.SetAPI(&api)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "!version",
		IsPrivate: false,
	}

	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   api.GetVersion(),
		IsPrivate: false,
	}
	p.OnCommand("version", "", postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}
