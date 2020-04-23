package voteplugin

import (
	"fmt"
	"strconv"
	"strings"
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
	storage := plugin.MockStorageAPI{}
	p.SetAPI(&api, &storage)
}

func TestVotePlugin_HelpTextAndInvalidCommands(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
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
	postToPlugin.Content = "!vote"
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   helpText,
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!votehelp"
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!vote    "
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!vote something"
	postToPlugin.IsPrivate = true
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)
}

func TestVotePlugin_FailOnPost(t *testing.T) {
	assert := assert.New(t)

	expectedChannel := "CHANNEL ID"
	expectedMessageID := "SOME MESSAGE ID"

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	storage := plugin.MockStorageAPI{}
	p.SetAPI(&api, &storage)
	api.ErrorToReturn = fmt.Errorf("Some error")

	api.PostResponse.PostedMessageIdent.Channel = expectedChannel
	api.PostResponse.PostedMessageIdent.ID = expectedMessageID

	postToPlugin := model.Post{
		ChannelID: expectedChannel,
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		IsPrivate: false,
	}

	api.Reset()
	voteText := "hello this is a vote"
	postToPlugin.Content = "!vote " + voteText
	expectedPostFromPlugin := model.Post{
		ChannelID: expectedChannel,
		Content:   "Sorry to inform you, but we failed to create the Vote! Please try again later.",
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
	storage := plugin.MockStorageAPI{}
	p.SetAPI(&api, &storage)

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
	expectedPostFromPlugin := model.Post{
		ChannelID: expectedChannel,
		Content:   voteEndHelpText,
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!voteend something else"
	expectedPostFromPlugin = model.Post{
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
	expectedPostFromPlugin = model.Post{
		ChannelID: expectedChannel,
		Content:   voteEndHelpText,
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

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

func TestVotePlugin_SimpleVoteCounting(t *testing.T) {
	assert := assert.New(t)

	expectedChannel := "CHANNEL ID"
	expectedMessageID := "SOME MESSAGE ID"

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	storage := plugin.MockStorageAPI{}
	p.SetAPI(&api, &storage)

	api.PostResponse.PostedMessageIdent.Channel = expectedChannel
	api.PostResponse.PostedMessageIdent.ID = expectedMessageID

	postToPlugin := model.Post{
		ChannelID: expectedChannel,
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		IsPrivate: false,
	}

	api.Reset()
	voteText := "hello this is a vote"
	postToPlugin.Content = "!vote " + voteText
	expectedPostFromPlugin := model.Post{
		ChannelID: expectedChannel,
		Content:   "\n*" + voteText + "*\n:one:: Yes\n:two:: No\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	expectedMessageIDFromPlugin := model.MessageIdentifier{
		ID:      expectedMessageID,
		Channel: expectedChannel,
	}
	reaction := model.Reaction{
		Message:  expectedMessageIDFromPlugin,
		Type:     "added",
		Reaction: "one",
		User:     model.User{Name: "USER 1"},
	}
	expectedPostFromPlugin = model.Post{
		ChannelID: expectedChannel,
		Content:   "\n*" + voteText + "*\n:one:: Yes - 1\n:two:: No\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!",
		IsPrivate: false,
	}
	p.OnReactionAdded(reaction)
	assert.Equal(true, api.WasUpdatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastUpdatePostPost)
	assert.Equal(expectedMessageIDFromPlugin, api.LastUpdatePostMessageID)

	api.Reset()
	reaction.Reaction = "two"
	expectedPostFromPlugin.Content = "\n*" + voteText + "*\n:one:: Yes - 1\n:two:: No - 1\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!"
	p.OnReactionAdded(reaction)
	assert.Equal(true, api.WasUpdatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastUpdatePostPost)
	assert.Equal(expectedMessageIDFromPlugin, api.LastUpdatePostMessageID)

	api.Reset()
	reaction.Type = "removed"
	reaction.Reaction = "one"
	expectedPostFromPlugin.Content = "\n*" + voteText + "*\n:one:: Yes\n:two:: No - 1\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!"
	p.OnReactionRemoved(reaction)
	assert.Equal(true, api.WasUpdatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastUpdatePostPost)
	assert.Equal(expectedMessageIDFromPlugin, api.LastUpdatePostMessageID)

	api.Reset()
	postToPlugin.Content = "!voteend " + voteText
	expectedPostFromPlugin.Content = "\n*" + voteText + "*\n:one:: Yes\n:two:: No - 1\nThis vote has ended, thanks for participating!"
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasUpdatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastUpdatePostPost)
	assert.Equal(expectedMessageIDFromPlugin, api.LastUpdatePostMessageID)
}

func TestVotePlugin_DoNotAllowCreationOfTheSameVoteTwice(t *testing.T) {
	assert := assert.New(t)

	expectedChannel := "CHANNEL ID"
	expectedOtherChannel := "OTHER CHANNEL ID"
	expectedMessageID := "SOME MESSAGE ID"
	expectedOtherMessageID := "SOME OTHER MESSAGE ID"

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	storage := plugin.MockStorageAPI{}
	p.SetAPI(&api, &storage)

	api.PostResponse.PostedMessageIdent.Channel = expectedChannel
	api.PostResponse.PostedMessageIdent.ID = expectedMessageID

	postToPlugin := model.Post{
		ChannelID: expectedChannel,
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		IsPrivate: false,
	}

	api.Reset()
	voteText := "hello this is a vote"
	postToPlugin.Content = "!vote " + voteText
	expectedPostFromPlugin := model.Post{
		ChannelID: expectedChannel,
		Content:   "\n*" + voteText + "*\n:one:: Yes\n:two:: No\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	voteText = "hello this is a vote"
	postToPlugin.Content = "!vote " + voteText
	expectedPostFromPlugin = model.Post{
		ChannelID: expectedChannel,
		Content:   "A vote with the same description is already running. End that vote first or enter a different description.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	voteText = "hello this is a vote"
	postToPlugin.Content = "!vote " + voteText
	postToPlugin.ChannelID = expectedOtherChannel
	api.PostResponse.PostedMessageIdent.Channel = expectedOtherChannel
	api.PostResponse.PostedMessageIdent.ID = expectedOtherMessageID
	expectedPostFromPlugin = model.Post{
		ChannelID: expectedOtherChannel,
		Content:   "\n*" + voteText + "*\n:one:: Yes\n:two:: No\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	// If the vote has ended, we shall be allowed to create it again
	api.Reset()
	postToPlugin.Content = "!voteend " + voteText
	postToPlugin.ChannelID = expectedChannel
	api.PostResponse.PostedMessageIdent.Channel = expectedChannel
	api.PostResponse.PostedMessageIdent.ID = expectedMessageID
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

	api.Reset()
	voteText = "hello this is a vote"
	postToPlugin.Content = "!vote " + voteText
	expectedPostFromPlugin = model.Post{
		ChannelID: expectedChannel,
		Content:   "\n*" + voteText + "*\n:one:: Yes\n:two:: No\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestVotePlugin_CreateAndEndCustomVote(t *testing.T) {
	assert := assert.New(t)

	expectedChannel := "CHANNEL ID"
	expectedMessageID := "SOME MESSAGE ID"
	customOptions := []string{"red", "green", "blue"}
	customOptionsStr := "[" + strings.Join(customOptions, ",") + "]"

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	storage := plugin.MockStorageAPI{}
	p.SetAPI(&api, &storage)

	api.PostResponse.PostedMessageIdent.Channel = expectedChannel
	api.PostResponse.PostedMessageIdent.ID = expectedMessageID

	postToPlugin := model.Post{
		ChannelID: expectedChannel,
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		IsPrivate: false,
	}

	api.Reset()
	voteText := "hello this is a vote"
	postToPlugin.Content = "!vote " + voteText + " " + customOptionsStr
	expectedPostFromPlugin := model.Post{
		ChannelID: expectedChannel,
		Content:   "\n*" + voteText + "*\n:one:: " + customOptions[0] + "\n:two:: " + customOptions[1] + "\n:three:: " + customOptions[2] + "\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!voteend " + voteText
	expectedPostFromPlugin = model.Post{
		ChannelID: expectedChannel,
		Content:   "\n*" + voteText + "*\n:one:: " + customOptions[0] + "\n:two:: " + customOptions[1] + "\n:three:: " + customOptions[2] + "\nThis vote has ended, thanks for participating!",
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

	// Empty string, i.e., ',' at the end of custom options shall not lead to another empty option
	customOptions = append(customOptions, "")
	customOptionsStr = "[" + strings.Join(customOptions, ",") + "]"

	api.Reset()
	voteText = "hello this is another vote"
	postToPlugin.Content = "!vote " + voteText + " " + customOptionsStr
	expectedPostFromPlugin = model.Post{
		ChannelID: expectedChannel,
		Content:   "\n*" + voteText + "*\n:one:: " + customOptions[0] + "\n:two:: " + customOptions[1] + "\n:three:: " + customOptions[2] + "\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	// Empty options, i.e., ',,' somewhere in the custom options shall not lead to another empty option
	customOptions = append(customOptions, "")
	customOptionsStr = "[" + customOptions[0] + ",," + customOptions[1] + "," + customOptions[2] + "]"

	api.Reset()
	voteText = "hello this is yet another vote"
	postToPlugin.Content = "!vote " + voteText + " " + customOptionsStr
	expectedPostFromPlugin = model.Post{
		ChannelID: expectedChannel,
		Content:   "\n*" + voteText + "*\n:one:: " + customOptions[0] + "\n:two:: " + customOptions[1] + "\n:three:: " + customOptions[2] + "\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestVotePlugin_CustomVoteCounting(t *testing.T) {
	assert := assert.New(t)

	expectedChannel := "CHANNEL ID"
	expectedMessageID := "SOME MESSAGE ID"
	customOptions := []string{"red", "green", "blue"}
	customOptionsStr := "[" + strings.Join(customOptions, ",") + "]"

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	storage := plugin.MockStorageAPI{}
	p.SetAPI(&api, &storage)

	api.PostResponse.PostedMessageIdent.Channel = expectedChannel
	api.PostResponse.PostedMessageIdent.ID = expectedMessageID

	postToPlugin := model.Post{
		ChannelID: expectedChannel,
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		IsPrivate: false,
	}

	api.Reset()
	voteText := "hello this is a vote"
	postToPlugin.Content = "!vote " + voteText + " " + customOptionsStr
	expectedPostFromPlugin := model.Post{
		ChannelID: expectedChannel,
		Content:   "\n*" + voteText + "*\n:one:: " + customOptions[0] + "\n:two:: " + customOptions[1] + "\n:three:: " + customOptions[2] + "\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	expectedMessageIDFromPlugin := model.MessageIdentifier{
		ID:      expectedMessageID,
		Channel: expectedChannel,
	}
	reaction := model.Reaction{
		Message:  expectedMessageIDFromPlugin,
		Type:     "added",
		Reaction: "one",
		User:     model.User{Name: "USER 1"},
	}
	expectedPostFromPlugin = model.Post{
		ChannelID: expectedChannel,
		Content:   "\n*" + voteText + "*\n:one:: " + customOptions[0] + " - 1\n:two:: " + customOptions[1] + "\n:three:: " + customOptions[2] + "\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!",
		IsPrivate: false,
	}
	p.OnReactionAdded(reaction)
	assert.Equal(true, api.WasUpdatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastUpdatePostPost)
	assert.Equal(expectedMessageIDFromPlugin, api.LastUpdatePostMessageID)

	api.Reset()
	reaction.Reaction = "three"
	expectedPostFromPlugin.Content = "\n*" + voteText + "*\n:one:: " + customOptions[0] + " - 1\n:two:: " + customOptions[1] + "\n:three:: " + customOptions[2] + " - 1\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!"
	p.OnReactionAdded(reaction)
	assert.Equal(true, api.WasUpdatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastUpdatePostPost)
	assert.Equal(expectedMessageIDFromPlugin, api.LastUpdatePostMessageID)

	api.Reset()
	reaction.Type = "removed"
	reaction.Reaction = "one"
	expectedPostFromPlugin.Content = "\n*" + voteText + "*\n:one:: " + customOptions[0] + "\n:two:: " + customOptions[1] + "\n:three:: " + customOptions[2] + " - 1\nParticipate by reacting with the appropriate emoji corresponding to the option you want to vote for!"
	p.OnReactionRemoved(reaction)
	assert.Equal(true, api.WasUpdatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastUpdatePostPost)
	assert.Equal(expectedMessageIDFromPlugin, api.LastUpdatePostMessageID)

	api.Reset()
	postToPlugin.Content = "!voteend " + voteText
	expectedPostFromPlugin.Content = "\n*" + voteText + "*\n:one:: " + customOptions[0] + "\n:two:: " + customOptions[1] + "\n:three:: " + customOptions[2] + " - 1\nThis vote has ended, thanks for participating!"
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasUpdatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastUpdatePostPost)
	assert.Equal(expectedMessageIDFromPlugin, api.LastUpdatePostMessageID)
}

func TestVotePlugin_CustomVoteOptionsLimit(t *testing.T) {
	assert := assert.New(t)

	optionsLimit := 10
	overLimit := optionsLimit + 1
	expectedChannel := "CHANNEL ID"
	expectedMessageID := "SOME MESSAGE ID"
	customOptionsText := "["
	for i := 0; i < overLimit; i++ {
		customOptionsText += strconv.Itoa(i)
		if i < overLimit-1 {
			customOptionsText += ","
		}
	}
	customOptionsText += "]"

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	storage := plugin.MockStorageAPI{}
	p.SetAPI(&api, &storage)

	api.PostResponse.PostedMessageIdent.Channel = expectedChannel
	api.PostResponse.PostedMessageIdent.ID = expectedMessageID

	postToPlugin := model.Post{
		ChannelID: expectedChannel,
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		IsPrivate: false,
	}

	api.Reset()
	voteText := "hello this is a vote"
	postToPlugin.Content = "!vote " + voteText + " " + customOptionsText
	expectedPostFromPlugin := model.Post{
		ChannelID: expectedChannel,
		Content:   "More than the allowed number of options specified. Please specify " + strconv.Itoa(optionsLimit) + " or less options.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}
