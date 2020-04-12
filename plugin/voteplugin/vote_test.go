package voteplugin

import (
	"testing"

	"git.abyle.org/redseligg/botorchestrator/botconfig"
	"github.com/stretchr/testify/assert"
	"github.com/torlenor/abylebotter/plugin"
)

func TestCreateVotePlugin(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: "something"})
	assert.Error(err)
	assert.Nil(p)

	p, err = New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	p.SetAPI(&api)
}

func TestVotePlugin_OnPost(t *testing.T) {
	// assert := assert.New(t)

	// p := VotePlugin{
	// 	randomizer: mockRandomizer{},
	// }
	// assert.Equal(nil, p.API)

	// api := plugin.MockAPI{}
	// p.SetAPI(&api)

	// postToPlugin := model.Post{
	// 	ChannelID: "CHANNEL ID",
	// 	Channel:   "SOME CHANNEL",
	// 	UserID:    "SOME USER ID",
	// 	User:      "USER 1",
	// 	Content:   "MESSAGE CONTENT",
	// 	IsPrivate: false,
	// }
	// p.OnPost(postToPlugin)
	// assert.Equal(false, api.WasCreatePostCalled)

	// api.Reset()
	// postToPlugin.Content = "!roll"
	// expectedPostFromPlugin := model.Post{
	// 	ChannelID: "CHANNEL ID",
	// 	Channel:   "SOME CHANNEL",
	// 	UserID:    "SOME USER ID",
	// 	User:      "USER 1",
	// 	Content:   "<@" + postToPlugin.UserID + "> rolled *" + strconv.Itoa(123) + "* in [0,100]",
	// 	IsPrivate: false,
	// }
	// p.OnPost(postToPlugin)
	// assert.Equal(true, api.WasCreatePostCalled)
	// assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	// api.Reset()
	// postToPlugin.Content = "!roll 1000"
	// expectedPostFromPlugin = model.Post{
	// 	ChannelID: "CHANNEL ID",
	// 	Channel:   "SOME CHANNEL",
	// 	UserID:    "SOME USER ID",
	// 	User:      "USER 1",
	// 	Content:   "<@" + postToPlugin.UserID + "> rolled *" + strconv.Itoa(123) + "* in [0,1000]",
	// 	IsPrivate: false,
	// }
	// p.OnPost(postToPlugin)
	// assert.Equal(true, api.WasCreatePostCalled)
	// assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	// api.Reset()
	// postToPlugin.Content = "!roll -1"
	// expectedPostFromPlugin = model.Post{
	// 	ChannelID: "CHANNEL ID",
	// 	Channel:   "SOME CHANNEL",
	// 	UserID:    "SOME USER ID",
	// 	User:      "USER 1",
	// 	Content:   "Number must be > 0",
	// 	IsPrivate: false,
	// }
	// p.OnPost(postToPlugin)
	// assert.Equal(true, api.WasCreatePostCalled)
	// assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	// api.Reset()
	// postToPlugin.Content = "!roll sdsadsad"
	// expectedPostFromPlugin = model.Post{
	// 	ChannelID: "CHANNEL ID",
	// 	Channel:   "SOME CHANNEL",
	// 	UserID:    "SOME USER ID",
	// 	User:      "USER 1",
	// 	Content:   "Not a number",
	// 	IsPrivate: false,
	// }
	// p.OnPost(postToPlugin)
	// assert.Equal(true, api.WasCreatePostCalled)
	// assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}
