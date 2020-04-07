package giveawayplugin

import (
	"testing"

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
		UserID:    "SOME USER ID",
		User:      "USER 1",
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
		Channel:   "",
		UserID:    "",
		User:      "",
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
		Channel:   "",
		UserID:    "",
		User:      "",
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
		Channel:   "",
		UserID:    "",
		User:      "",
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
		Channel:   "",
		UserID:    "",
		User:      "",
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
		Channel:   "",
		UserID:    "",
		User:      "",
		Content:   "Type '!gstart <time> <secretword> [winners] [prize]' to start a new giveaway.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
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
		UserID:    "SOME USER ID",
		User:      "USER 1",
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}

	api.Reset()
	postToPlugin.Content = "!gend"
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "",
		UserID:    "",
		User:      "",
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
		Channel:   "",
		UserID:    "",
		User:      "",
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
		Channel:   "",
		UserID:    "",
		User:      "",
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
		Channel:   "",
		UserID:    "",
		User:      "",
		Content:   "Giveaway started! Type " + secretword + " to participate.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	userPostToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		UserID:    "PARTICIPANT_1_ID",
		User:      "PARTICIPANT_1",
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
		Channel:   "",
		UserID:    "",
		User:      "",
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
		Channel:   "",
		UserID:    "",
		User:      "",
		Content:   "The winner(s) is/are <@" + userPostToPlugin.UserID + ">. Congratulations!",
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
		UserID:    "SOME USER ID",
		User:      "USER 1",
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}

	api.Reset()
	secretword := "sonne"
	postToPlugin.Content = "!gstart 5m " + secretword + " 2"
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "",
		UserID:    "",
		User:      "",
		Content:   "Giveaway started! Type " + secretword + " to participate.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	userPostToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		UserID:    "PARTICIPANT_1_ID",
		User:      "PARTICIPANT_1",
		Content:   secretword,
		IsPrivate: false,
	}

	api.Reset()
	p.OnPost(userPostToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	userPostToPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		UserID:    "PARTICIPANT_2_ID",
		User:      "PARTICIPANT_2",
		Content:   secretword,
		IsPrivate: false,
	}

	api.Reset()
	p.OnPost(userPostToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	userPostToPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		UserID:    "PARTICIPANT_3_ID",
		User:      "PARTICIPANT_3",
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
		Channel:   "",
		UserID:    "",
		User:      "",
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
		UserID:    "SOME USER ID",
		User:      "USER 1",
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}

	api.Reset()
	secretword := "sonne"
	prize := "That awesome PRIZE"
	postToPlugin.Content = "!gstart 5m " + secretword + " 1 " + prize
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "",
		UserID:    "",
		User:      "",
		Content:   "Giveaway started! Type " + secretword + " to participate.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	userPostToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		UserID:    "PARTICIPANT_1_ID",
		User:      "PARTICIPANT_1",
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
		Channel:   "",
		UserID:    "",
		User:      "",
		Content:   "The winner(s) is/are <@" + userPostToPlugin.UserID + ">. You won 'That awesome PRIZE'. Congratulations!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
	assert.Equal(randomizer.Argument, 1)
}

func TestGiveawayPluginCreateAndEndGiveawayAndReroll(t *testing.T) {
	// assert := assert.New(t)

	// p, err := New(botconfig.PluginConfig{Type: "giveaway"})
	// assert.NoError(err)
	// assert.Equal(nil, p.API)

	// randomizer := &mockRandomizer{}
	// p.randomizer = randomizer

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

	// api.Reset()
	// secretword := "sonne"
	// postToPlugin.Content = "!gstart 5m " + secretword + " 1"
	// expectedPostFromPlugin := model.Post{
	// 	ChannelID: "CHANNEL ID",
	// 	Channel:   "",
	// 	UserID:    "",
	// 	User:      "",
	// 	Content:   "Giveaway started! Type " + secretword + " to participate.",
	// 	IsPrivate: false,
	// }
	// p.OnPost(postToPlugin)
	// assert.Equal(true, api.WasCreatePostCalled)
	// assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	// userPostToPlugin := model.Post{
	// 	ChannelID: "CHANNEL ID",
	// 	Channel:   "SOME CHANNEL",
	// 	UserID:    "PARTICIPANT_1_ID",
	// 	User:      "PARTICIPANT_1",
	// 	Content:   secretword,
	// 	IsPrivate: false,
	// }

	// api.Reset()
	// p.OnPost(userPostToPlugin)
	// assert.Equal(false, api.WasCreatePostCalled)

	// userPostToPlugin = model.Post{
	// 	ChannelID: "CHANNEL ID",
	// 	Channel:   "SOME CHANNEL",
	// 	UserID:    "PARTICIPANT_2_ID",
	// 	User:      "PARTICIPANT_2",
	// 	Content:   secretword,
	// 	IsPrivate: false,
	// }

	// api.Reset()
	// p.OnPost(userPostToPlugin)
	// assert.Equal(false, api.WasCreatePostCalled)

	// userPostToPlugin = model.Post{
	// 	ChannelID: "CHANNEL ID",
	// 	Channel:   "SOME CHANNEL",
	// 	UserID:    "PARTICIPANT_3_ID",
	// 	User:      "PARTICIPANT_3",
	// 	Content:   secretword,
	// 	IsPrivate: false,
	// }

	// api.Reset()
	// p.OnPost(userPostToPlugin)
	// assert.Equal(false, api.WasCreatePostCalled)

	// api.Reset()
	// postToPlugin.Content = "!greroll"
	// expectedPostFromPlugin = model.Post{
	// 	ChannelID: "CHANNEL ID",
	// 	Channel:   "",
	// 	UserID:    "",
	// 	User:      "",
	// 	Content:   "Cannot reroll when there is still a giveaway ongoing.",
	// 	IsPrivate: false,
	// }
	// p.OnPost(postToPlugin)
	// assert.Equal(true, api.WasCreatePostCalled)
	// assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
	// assert.Equal(randomizer.Argument, 3)

	// api.Reset()
	// postToPlugin.Content = "!gend"
	// expectedPostFromPlugin = model.Post{
	// 	ChannelID: "CHANNEL ID",
	// 	Channel:   "",
	// 	UserID:    "",
	// 	User:      "",
	// 	Content:   "The winner(s) is/are <@" + "PARTICIPANT_1_ID" + ">. Congratulations!",
	// 	IsPrivate: false,
	// }
	// p.OnPost(postToPlugin)
	// assert.Equal(true, api.WasCreatePostCalled)
	// assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
	// assert.Equal(randomizer.Argument, 3)

	// api.Reset()
	// postToPlugin.Content = "!greroll"
	// expectedPostFromPlugin = model.Post{
	// 	ChannelID: "CHANNEL ID",
	// 	Channel:   "",
	// 	UserID:    "",
	// 	User:      "",
	// 	Content:   "The winner(s) is/are <@" + "PARTICIPANT_1_ID" + ">. Congratulations!",
	// 	IsPrivate: false,
	// }
	// p.OnPost(postToPlugin)
	// assert.Equal(true, api.WasCreatePostCalled)
	// assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
	// assert.Equal(randomizer.Argument, 3)
}
