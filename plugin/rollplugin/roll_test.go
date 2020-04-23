package rollplugin

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/plugin"
)

type mockRandomizer struct{}

func (r mockRandomizer) random(int) int {
	return 123
}

func TestCreateRollPlugin(t *testing.T) {
	assert := assert.New(t)

	p, err := New()
	assert.NoError(err)
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	storage := plugin.MockStorageAPI{}
	p.SetAPI(&api, &storage)
}

func TestRollPlugin_OnPost(t *testing.T) {
	assert := assert.New(t)

	p := RollPlugin{
		randomizer: mockRandomizer{},
	}
	assert.Equal(nil, p.API)

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
	postToPlugin.Content = "!roll"
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "<@" + postToPlugin.User.ID + "> rolled *" + strconv.Itoa(123) + "* in [0,100]",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!roll 1000"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "<@" + postToPlugin.User.ID + "> rolled *" + strconv.Itoa(123) + "* in [0,1000]",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!roll -1"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "Number must be > 0",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!roll sdsadsad"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "Not a number",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}
