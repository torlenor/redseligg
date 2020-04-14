package giveawayplugin

import (
	"testing"
	"time"

	"git.abyle.org/redseligg/botorchestrator/botconfig"

	"github.com/stretchr/testify/assert"
	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/plugin"
)

type mockRandomizer struct {
	RandomNumber int

	Argument int
}

func (r *mockRandomizer) Intn(arg int) int {
	r.Argument = arg
	return r.RandomNumber
}

func (r *mockRandomizer) Shuffle(arg int, swap func(i, j int)) {
	r.Argument = arg
}

func TestCreateGiveawayPlugin(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: "something"})
	assert.Error(err)
	assert.Nil(p)

	p, err = New(botconfig.PluginConfig{Type: "giveaway"})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	p.SetAPI(&api)
}

func TestGiveawayPluginHelpTextAndInvalidCommands(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: "giveaway"})
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
	postToPlugin.Content = "!gstarttt"
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.Content = "!gstart help"
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Type '!gstart <time> <secretword> [winners] [prize]' to start a new giveaway.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!gstart    "
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Type '!gstart <time> <secretword> [winners] [prize]' to start a new giveaway.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!gstart 1m"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Type '!gstart <time> <secretword> [winners] [prize]' to start a new giveaway.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!gstart 1kk hello"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Type '!gstart <time> <secretword> [winners] [prize]' to start a new giveaway.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!gstart 1m hello k"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Type '!gstart <time> <secretword> [winners] [prize]' to start a new giveaway.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!gstart 1m hello"
	postToPlugin.IsPrivate = true
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)
}

func TestGiveawayPluginCreateAndEndGiveaway(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: "giveaway"})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	randomizer := &mockRandomizer{}
	p.randomizer = randomizer

	api := plugin.MockAPI{}
	p.SetAPI(&api)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}

	api.Reset()
	postToPlugin.Content = "!gend"
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "No giveaway running. Use !gstart command to start a new one.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	secretword := "hello"
	postToPlugin.Content = "!gstart 10m " + secretword
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Giveaway started! Type " + secretword + " to participate.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!gend"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Cannot pick a winner. There were no participants to the giveaway.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	secretword = "sonne"
	postToPlugin.Content = "!gstart 5m " + secretword
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Giveaway started! Type " + secretword + " to participate.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	userPostToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "PARTICIPANT_1_ID", Name: "PARTICIPANT_1"},
		Content:   secretword,
		IsPrivate: false,
	}

	api.Reset()
	p.OnPost(userPostToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	secretword = "hello"
	postToPlugin.Content = "!gstart 10m " + secretword
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Giveaway already running.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!gend"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "The winner(s) is/are <@" + userPostToPlugin.User.ID + ">. Congratulations!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
	assert.Equal(randomizer.Argument, 1)
}

func TestGiveawayPluginCreateAndAutomaticEndGiveaway(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: "giveaway"})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	randomizer := &mockRandomizer{}
	p.randomizer = randomizer

	api := plugin.MockAPI{}
	p.SetAPI(&api)
	p.OnRun()

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}

	api.Reset()
	secretword := "sonne"
	postToPlugin.Content = "!gstart 2s " + secretword
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Giveaway started! Type " + secretword + " to participate.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	userPostToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "PARTICIPANT_1_ID", Name: "PARTICIPANT_1"},
		Content:   secretword,
		IsPrivate: false,
	}

	api.Reset()
	p.OnPost(userPostToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	// TODO refactor that to not use that ugly sleep here
	time.Sleep(4 * time.Second)
	assert.Equal(true, api.WasCreatePostCalled)
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "The winner(s) is/are <@" + userPostToPlugin.User.ID + ">. Congratulations!",
		IsPrivate: false,
	}
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
	assert.Equal(randomizer.Argument, 1)

	p.OnStop()
	// TODO refactor that to not use that ugly sleep here
	time.Sleep(1 * time.Second)
}

func TestGiveawayPluginCreateAndEndGiveawayOnlyMods(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: "giveaway"})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	allowedUser := "User 2"
	notAllowedUser := "User 1"

	p.cfg.OnlyMods = true
	p.cfg.Mods = []string{allowedUser}

	randomizer := &mockRandomizer{}
	p.randomizer = randomizer

	api := plugin.MockAPI{}
	p.SetAPI(&api)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: notAllowedUser},
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}

	api.Reset()
	secretword := "hello"
	postToPlugin.Content = "!gstart 10m " + secretword
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	secretword = "hello"
	postToPlugin.User = model.User{ID: "SOME USER ID", Name: allowedUser}
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Giveaway started! Type " + secretword + " to participate.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	userPostToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "PARTICIPANT_1_ID", Name: "PARTICIPANT_1"},
		Content:   secretword,
		IsPrivate: false,
	}

	api.Reset()
	p.OnPost(userPostToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.Content = "!gend"
	postToPlugin.User = model.User{ID: "SOME USER ID", Name: notAllowedUser}
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.User = model.User{ID: "SOME USER ID", Name: allowedUser}
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "The winner(s) is/are <@" + userPostToPlugin.User.ID + ">. Congratulations!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
	assert.Equal(randomizer.Argument, 1)
}

func TestGiveawayPluginCreateAndEndGiveawayWithMultipleWinners(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: "giveaway"})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	randomizer := &mockRandomizer{}
	p.randomizer = randomizer

	api := plugin.MockAPI{}
	p.SetAPI(&api)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}

	api.Reset()
	secretword := "sonne"
	postToPlugin.Content = "!gstart 5m " + secretword + " 2"
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Giveaway started! Type " + secretword + " to participate.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	userPostToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "PARTICIPANT_1_ID", Name: "PARTICIPANT_1"},
		Content:   secretword,
		IsPrivate: false,
	}

	api.Reset()
	p.OnPost(userPostToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	userPostToPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "PARTICIPANT_2_ID", Name: "PARTICIPANT_2"},
		Content:   secretword,
		IsPrivate: false,
	}

	api.Reset()
	p.OnPost(userPostToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	userPostToPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "PARTICIPANT_3_ID", Name: "PARTICIPANT_3"},
		Content:   secretword,
		IsPrivate: false,
	}

	api.Reset()
	p.OnPost(userPostToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.Content = "!gend"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "The winner(s) is/are <@" + "PARTICIPANT_1_ID" + ">, <@" + "PARTICIPANT_2_ID" + ">. Congratulations!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
	assert.Equal(randomizer.Argument, 3)
}

func TestGiveawayPluginCreateAndEndGiveawayWithPrize(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: "giveaway"})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	randomizer := &mockRandomizer{}
	p.randomizer = randomizer

	api := plugin.MockAPI{}
	p.SetAPI(&api)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}

	api.Reset()
	secretword := "sonne"
	prize := "That awesome PRIZE"
	postToPlugin.Content = "!gstart 5m " + secretword + " 1 " + prize
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Giveaway started! Type " + secretword + " to participate.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	userPostToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "PARTICIPANT_1_ID", Name: "PARTICIPANT_1"},
		Content:   secretword,
		IsPrivate: false,
	}

	api.Reset()
	p.OnPost(userPostToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.Content = "!gend"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "The winner(s) is/are <@" + userPostToPlugin.User.ID + ">. You won 'That awesome PRIZE'. Congratulations!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
	assert.Equal(randomizer.Argument, 1)
}

func TestGiveawayPluginCreateAndEndGiveawayAndReroll(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: "giveaway"})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	randomizer := &mockRandomizer{}
	p.randomizer = randomizer

	api := plugin.MockAPI{}
	p.SetAPI(&api)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}

	api.Reset()
	postToPlugin.Content = "!greroll"
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "No previous giveaway in that channel. Use !gstart command to start a new one.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	secretword := "hello"
	postToPlugin.Content = "!gstart 10m " + secretword
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Giveaway started! Type " + secretword + " to participate.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!greroll"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Cannot pick a new winner. There is currently a giveaway running in this channel.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	userPostToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "PARTICIPANT_1_ID", Name: "PARTICIPANT_1"},
		Content:   secretword,
		IsPrivate: false,
	}

	api.Reset()
	p.OnPost(userPostToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.Content = "!gend"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "The winner(s) is/are <@" + userPostToPlugin.User.ID + ">. Congratulations!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
	assert.Equal(randomizer.Argument, 1)

	api.Reset()
	postToPlugin.Content = "!greroll"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "The new winner is <@" + userPostToPlugin.User.ID + ">. Congratulations!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
	assert.Equal(randomizer.Argument, 1)

	api.Reset()
	prize := "SOME AWESOME PRIZE"
	postToPlugin.Content = "!gstart 10m " + secretword + " 1 " + prize
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Giveaway started! Type " + secretword + " to participate.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	userPostToPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "PARTICIPANT_1_ID", Name: "PARTICIPANT_1"},
		Content:   secretword,
		IsPrivate: false,
	}

	api.Reset()
	p.OnPost(userPostToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.Content = "!gend"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "The winner(s) is/are <@" + userPostToPlugin.User.ID + ">. You won '" + prize + "'. Congratulations!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
	assert.Equal(randomizer.Argument, 1)

	api.Reset()
	postToPlugin.Content = "!greroll"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "The new winner is <@" + userPostToPlugin.User.ID + ">. You won '" + prize + "'. Congratulations!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
	assert.Equal(randomizer.Argument, 1)

	api.Reset()
	secretword = "hello"
	postToPlugin.Content = "!gstart 10m " + secretword
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Giveaway started! Type " + secretword + " to participate.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!gend"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Cannot pick a winner. There were no participants to the giveaway.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!greroll"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Cannot pick a new winner. There were no participants to the previous giveaway.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}