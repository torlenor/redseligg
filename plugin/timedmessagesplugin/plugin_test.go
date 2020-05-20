package timedmessagesplugin

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/torlenor/redseligg/botconfig"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/platform"
	"github.com/torlenor/redseligg/plugin"
	"github.com/torlenor/redseligg/storagemodels"
)

var providedFeatures = map[string]bool{
	platform.FeatureMessagePost: true,
}

var command = "tm"

func TestCreateTimedMessagesPlugin(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: "something"})
	assert.Error(err)
	assert.Nil(p)

	p, err = New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.NotNil(p)
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	err = p.SetAPI(&api)
	assert.Error(err)

	api.ProvidedFeatures = providedFeatures
	err = p.SetAPI(&api)
	assert.NoError(err)

	assert.Equal(PLUGIN_TYPE, p.PluginType())
}

func TestTimedMessagesPlugin_HasExpectedRequiredFeatures(t *testing.T) {
	assert := assert.New(t)

	expectedRequiredFeatures := []string{
		platform.FeatureMessagePost,
	}

	p, _ := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.Equal(expectedRequiredFeatures, p.NeededFeatures)
}

func TestTimedMessagesPlugin_OnRun(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	storage := &MockStorage{}
	api := plugin.MockAPI{Storage: nil}
	api.ProvidedFeatures = providedFeatures
	err = p.SetAPI(&api)
	assert.NoError(err)

	assert.Equal("", api.LastLoggedError)
	p.OnRun()

	assert.Equal(ErrNoValidStorage.Error(), api.LastLoggedError)

	api.Reset()
	api.Storage = storage
	p.OnRun()
	assert.Equal("", api.LastLoggedError)
}

func TestTimedMessagesPlugin_HelpTextAndInvalidCommands(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	storage := &MockStorage{}
	api := plugin.MockAPI{Storage: storage}
	api.ProvidedFeatures = providedFeatures
	err = p.SetAPI(&api)
	assert.NoError(err)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "", // not used by plugin
		IsPrivate: false,
	}

	api.Reset()
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   helpText,
		IsPrivate: false,
	}
	p.OnCommand(command, "", postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	expectedPostFromPlugin.Content = helpTextAdd
	p.OnCommand(command, "add", postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	expectedPostFromPlugin.Content = helpTextRemove
	p.OnCommand(command, "remove", postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestTimedMessagesPlugin_AddAndRemoveTimedMessage(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	storage := &MockStorage{}
	api := plugin.MockAPI{Storage: storage}
	api.ProvidedFeatures = providedFeatures
	err = p.SetAPI(&api)
	assert.NoError(err)

	timeIntervalStr := "1m"
	timeInterval, _ := time.ParseDuration(timeIntervalStr)

	message := "some message"

	api.Reset()
	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "", // not used by plugin
		IsPrivate: false,
	}
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   fmt.Sprintf("Timed message '%s' with interval %s added.", message, timeInterval),
		IsPrivate: false,
	}
	p.OnCommand(command, "add "+timeIntervalStr+" "+message, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	assert.Equal(p.BotID, storage.LastRetrieved.BotID)
	assert.Equal(p.PluginID, storage.LastRetrieved.PluginID)
	assert.Equal(identFieldTimedMessages, storage.LastRetrieved.Identifier)

	assert.Equal(1, len(storage.StoredMessages.Data.Messages))

	storage.Reset()

	otherTimeIntervalStr := "2m"
	otherMessage := "some other message"

	expectedPostFromPlugin.Content = "Timed message to remove does not exist."
	p.OnCommand(command, "remove "+otherTimeIntervalStr+" "+otherMessage, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	expectedPostFromPlugin.Content = "Timed message to remove does not exist."
	p.OnCommand(command, "remove "+otherTimeIntervalStr+" "+message, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	storage.DataToReturn.Messages = append(storage.DataToReturn.Messages, storagemodels.TimedMessagesPluginMessage{
		Text:      message,
		Interval:  timeInterval,
		ChannelID: postToPlugin.ChannelID,
	})

	storage.DataToReturn.Messages = append(storage.DataToReturn.Messages, storagemodels.TimedMessagesPluginMessage{
		Text:      "something else which should not be removed",
		Interval:  timeInterval,
		ChannelID: postToPlugin.ChannelID,
	})

	expectedPostFromPlugin.Content = fmt.Sprintf("Timed message '%s' with interval %s removed.", message, timeInterval)
	p.OnCommand(command, "remove "+timeIntervalStr+" "+message, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	assert.Equal(p.BotID, storage.LastRetrieved.BotID)
	assert.Equal(p.PluginID, storage.LastRetrieved.PluginID)
	assert.Equal(identFieldTimedMessages, storage.LastRetrieved.Identifier)

	assert.Equal(1, len(storage.StoredMessages.Data.Messages))
	assert.Equal(storage.DataToReturn.Messages[1], storage.StoredMessages.Data.Messages[0])

	storage.ErrorToReturn = errors.New("Some error")
	expectedPostFromPlugin.Content = "Could not add timed message. Please try again later."
	p.OnCommand(command, "add "+timeIntervalStr+" "+message, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	expectedPostFromPlugin.Content = "Could not remove timed message. Please try again later."
	p.OnCommand(command, "remove "+timeIntervalStr+" "+message, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestTimedMessagesPlugin_AddAndRemoveTimedMessage_OnlyMods(t *testing.T) {
	assert := assert.New(t)

	userName := "SOME USER NAME"

	pluginID := "SOME_PLUGIN_ID"
	botID := "SOME BOT ID"
	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	p.PluginID = pluginID
	p.BotID = botID

	p.cfg.OnlyMods = true
	p.cfg.Mods = append(p.cfg.Mods, userName)

	storage := &MockStorage{}
	api := plugin.MockAPI{Storage: storage}
	api.ProvidedFeatures = providedFeatures
	err = p.SetAPI(&api)
	assert.NoError(err)

	timeIntervalStr := "1m"
	timeInterval, _ := time.ParseDuration(timeIntervalStr)

	message := "some message"

	api.Reset()
	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "", // not used by plugin
		IsPrivate: false,
	}
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   fmt.Sprintf("Timed message '%s' with interval %s added.", message, timeInterval),
		IsPrivate: false,
	}
	content := "add " + timeIntervalStr + " " + message
	p.OnCommand(command, content, postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.User.Name = userName
	p.OnCommand(command, content, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	assert.Equal(p.BotID, storage.LastRetrieved.BotID)
	assert.Equal(p.PluginID, storage.LastRetrieved.PluginID)
	assert.Equal(identFieldTimedMessages, storage.LastRetrieved.Identifier)

	assert.Equal(1, len(storage.StoredMessages.Data.Messages))

	storage.Reset()

	storage.DataToReturn.Messages = append(storage.DataToReturn.Messages, storagemodels.TimedMessagesPluginMessage{
		Text:      message,
		Interval:  timeInterval,
		ChannelID: postToPlugin.ChannelID,
	})

	storage.DataToReturn.Messages = append(storage.DataToReturn.Messages, storagemodels.TimedMessagesPluginMessage{
		Text:      "something else which should not be removed",
		Interval:  timeInterval,
		ChannelID: postToPlugin.ChannelID,
	})

	api.Reset()
	postToPlugin.User.Name = " some other user"
	content = "remove " + timeIntervalStr + " " + message
	expectedPostFromPlugin.Content = fmt.Sprintf("Timed message '%s' with interval %s removed.", message, timeInterval)

	p.OnCommand(command, content, postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.User.Name = userName
	p.OnCommand(command, content, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	assert.Equal(p.BotID, storage.LastRetrieved.BotID)
	assert.Equal(p.PluginID, storage.LastRetrieved.PluginID)
	assert.Equal(identFieldTimedMessages, storage.LastRetrieved.Identifier)

	assert.Equal(1, len(storage.StoredMessages.Data.Messages))
	assert.Equal(storage.DataToReturn.Messages[1], storage.StoredMessages.Data.Messages[0])

	storage.ErrorToReturn = errors.New("Some error")
	content = "add " + timeIntervalStr + " " + message
	expectedPostFromPlugin.Content = "Could not add timed message. Please try again later."
	p.OnCommand(command, content, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	content = "remove " + timeIntervalStr + " " + message
	expectedPostFromPlugin.Content = "Could not remove timed message. Please try again later."
	p.OnCommand(command, content, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestTimedMessagesPlugin_AddAndRemoveTimedMessage_AllMessages(t *testing.T) {
	assert := assert.New(t)

	message := "some message"

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	storage := &MockStorage{}
	api := plugin.MockAPI{Storage: storage}
	api.ProvidedFeatures = providedFeatures
	err = p.SetAPI(&api)
	assert.NoError(err)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "", // not used by plugin
		IsPrivate: false,
	}

	content := "remove all " + message

	storage.DataToReturn.Messages = append(storage.DataToReturn.Messages, storagemodels.TimedMessagesPluginMessage{
		Text:      message,
		Interval:  time.Duration(5) * time.Second,
		ChannelID: postToPlugin.ChannelID,
	})

	storage.DataToReturn.Messages = append(storage.DataToReturn.Messages, storagemodels.TimedMessagesPluginMessage{
		Text:      message,
		Interval:  time.Duration(15) * time.Second,
		ChannelID: postToPlugin.ChannelID,
	})

	storage.DataToReturn.Messages = append(storage.DataToReturn.Messages, storagemodels.TimedMessagesPluginMessage{
		Text:      "something else which should not be removed",
		Interval:  time.Duration(10) * time.Second,
		ChannelID: postToPlugin.ChannelID,
	})

	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   fmt.Sprintf("All timed messages with text '%s' removed.", message),
		IsPrivate: false,
	}
	p.OnCommand(command, content, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	assert.Equal(p.BotID, storage.LastRetrieved.BotID)
	assert.Equal(p.PluginID, storage.LastRetrieved.PluginID)
	assert.Equal(identFieldTimedMessages, storage.LastRetrieved.Identifier)

	assert.Equal(1, len(storage.StoredMessages.Data.Messages))
	assert.Equal(storage.DataToReturn.Messages[2], storage.StoredMessages.Data.Messages[0])
}

func TestTimedMessagesPlugin_NoStorage(t *testing.T) {
	assert := assert.New(t)

	pluginID := "SOME_PLUGIN_ID"
	botID := "SOME BOT ID"
	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	p.PluginID = pluginID
	p.BotID = botID

	api := plugin.MockAPI{Storage: nil}
	api.ProvidedFeatures = providedFeatures
	err = p.SetAPI(&api)
	assert.NoError(err)

	timeIntervalStr := "1m"

	message := "some message"

	api.Reset()
	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "", // not used by plugin
		IsPrivate: false,
	}
	content := "add " + timeIntervalStr + " " + message
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Could not add timed message. Please try again later.",
		IsPrivate: false,
	}
	p.OnCommand(command, content, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}
