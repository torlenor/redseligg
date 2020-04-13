package voteplugin

import (
	"testing"

	"git.abyle.org/redseligg/botorchestrator/botconfig"
	"github.com/stretchr/testify/assert"
	"github.com/torlenor/abylebotter/model"
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

func TestVotePlugin_HelpTextAndInvalidCommands(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
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
	postToPlugin.Content = "!vote"
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   HELP_TEXT,
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!votehelp"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   HELP_TEXT,
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!vote    "
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   HELP_TEXT,
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestVotePlugin_CreateAndEndSimpleVote(t *testing.T) {
	assert := assert.New(t)

	expectedChannel := "CHANNEL ID"
	expectedMessageID := "SOME MESSAGE ID"

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	p.SetAPI(&api)

	api.PostResponse.PostedMessageIdent.Channel = expectedChannel
	api.PostResponse.PostedMessageIdent.ID = expectedMessageID

	postToPlugin := model.Post{
		ChannelID: expectedChannel,
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		IsPrivate: false,
	}

	api.Reset()
	postToPlugin.Content = "!voteend"
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.Content = "!voteend something else"
	expectedPostFromPlugin := model.Post{
		ChannelID: expectedChannel,
		Content:   "No vote running with that description in this channel. Use the !vote command to start a new one.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	voteText := "hello this is a vote"
	postToPlugin.Content = "!vote " + voteText
	expectedPostFromPlugin = model.Post{
		ChannelID: expectedChannel,
		Content:   "\n*" + voteText + "*\n:one:: Yes\n:two:: No\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!voteend"
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.Content = "!voteend " + voteText
	expectedPostFromPlugin = model.Post{
		ChannelID: expectedChannel,
		Content:   "\n*" + voteText + "*\n:one:: Yes\n:two:: No\nThis vote has ended, thanks for participating!",
		IsPrivate: false,
	}
	expectedMessageIDFromPlugin := model.MessageIdentifier{
		ID:      expectedMessageID,
		Channel: expectedChannel,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasUpdatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastUpdatePostPost)
	assert.Equal(expectedMessageIDFromPlugin, api.LastUpdatePostMessageID)
}
