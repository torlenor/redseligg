package customcommandsplugin

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/torlenor/abylebotter/botconfig"

	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/platform"
	"github.com/torlenor/abylebotter/plugin"
	"github.com/torlenor/abylebotter/storagemodels"
)

var providedFeatures = map[string]bool{
	platform.FeatureMessagePost: true,
}

const command = "!customcommand"

func TestCreateCustomCommandsPlugin(t *testing.T) {
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

func TestCustomCommandsPlugin_HasExpectedRequiredFeatures(t *testing.T) {
	assert := assert.New(t)

	expectedRequiredFeatures := []string{
		platform.FeatureMessagePost,
	}

	p, _ := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.Equal(expectedRequiredFeatures, p.NeededFeatures)
}

func TestCustomCommandsPlugin_OnRun(t *testing.T) {
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

func TestCustomCommandsPlugin_HelpTextAndInvalidCommands(t *testing.T) {
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
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.Content = command
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   helpText,
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!customcommand add"
	expectedPostFromPlugin.Content = helpTextAdd
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!customcommand remove"
	expectedPostFromPlugin.Content = helpTextRemove
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestCustomCommandsPlugin_AddAndRemoveCustomCommand(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	storage := &MockStorage{}
	api := plugin.MockAPI{Storage: storage}
	api.ProvidedFeatures = providedFeatures
	err = p.SetAPI(&api)
	assert.NoError(err)

	customCommand := "someCustomCommand"
	message := "some message"

	api.Reset()
	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "!customcommand add " + customCommand + " " + message,
		IsPrivate: false,
	}
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   fmt.Sprintf("Custom command '%s' with message '%s' added.", customCommand, message),
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	assert.Equal(p.BotID, storage.LastRetrieved.BotID)
	assert.Equal(p.PluginID, storage.LastRetrieved.PluginID)
	assert.Equal(identField, storage.LastRetrieved.Identifier)

	assert.Equal(1, len(storage.StoredData.Data.Commands))

	storage.Reset()

	otherCustomCommand := "someOtherCommand"

	postToPlugin.Content = "!customcommand remove " + otherCustomCommand
	expectedPostFromPlugin.Content = "Custom command to remove does not exist."
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	storage.DataToReturn.Commands = append(storage.DataToReturn.Commands, storagemodels.CustomCommandsPluginCommand{
		Command:   customCommand,
		Text:      message,
		ChannelID: postToPlugin.ChannelID,
	})

	storage.DataToReturn.Commands = append(storage.DataToReturn.Commands, storagemodels.CustomCommandsPluginCommand{
		Command:   "commandNotToRemove",
		Text:      message,
		ChannelID: postToPlugin.ChannelID,
	})

	postToPlugin.Content = "!customcommand remove " + customCommand
	expectedPostFromPlugin.Content = fmt.Sprintf("Custom command '%s' removed.", customCommand)
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	assert.Equal(p.BotID, storage.LastRetrieved.BotID)
	assert.Equal(p.PluginID, storage.LastRetrieved.PluginID)
	assert.Equal(identField, storage.LastRetrieved.Identifier)

	assert.Equal(1, len(storage.StoredData.Data.Commands))
	assert.Equal(storage.DataToReturn.Commands[1], storage.StoredData.Data.Commands[0])

	storage.ErrorToReturn = errors.New("Some error")
	postToPlugin.Content = "!customcommand add " + customCommand + " " + message
	expectedPostFromPlugin.Content = "Could not add custom command. Please try again later."
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	postToPlugin.Content = "!customcommand remove " + customCommand
	expectedPostFromPlugin.Content = "Could not remove custom command. Please try again later."
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestCustomCommandsPlugin_UpdateCustomCommand(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	storage := &MockStorage{}
	api := plugin.MockAPI{Storage: storage}
	api.ProvidedFeatures = providedFeatures
	err = p.SetAPI(&api)
	assert.NoError(err)

	customCommand := "someCustomCommand"
	message := "some message"
	updatedMessage := "some other message"

	api.Reset()
	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "!customcommand add " + customCommand + " " + message,
		IsPrivate: false,
	}

	storage.DataToReturn.Commands = append(storage.DataToReturn.Commands, storagemodels.CustomCommandsPluginCommand{
		Command:   customCommand,
		Text:      message,
		ChannelID: postToPlugin.ChannelID,
	})

	api.Reset()
	postToPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "!customcommand add " + customCommand + " " + updatedMessage,
		IsPrivate: false,
	}
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   fmt.Sprintf("Custom command '%s' with message '%s' added.", customCommand, updatedMessage),
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	assert.Equal(p.BotID, storage.LastRetrieved.BotID)
	assert.Equal(p.PluginID, storage.LastRetrieved.PluginID)
	assert.Equal(identField, storage.LastRetrieved.Identifier)

	assert.Equal(1, len(storage.StoredData.Data.Commands))
	assert.Equal(customCommand, storage.StoredData.Data.Commands[0].Command)
	assert.Equal(updatedMessage, storage.StoredData.Data.Commands[0].Text)
}

func TestCustomCommandsPlugin_AddAndRemoveCustomCommand_OnlyMods(t *testing.T) {
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

	customCommand := "someCustomCommand"
	message := "some message"

	api.Reset()
	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "!customcommand add " + customCommand + " " + message,
		IsPrivate: false,
	}
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   fmt.Sprintf("Custom command '%s' with message '%s' added.", customCommand, message),
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.User.Name = userName
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	assert.Equal(p.BotID, storage.LastRetrieved.BotID)
	assert.Equal(p.PluginID, storage.LastRetrieved.PluginID)
	assert.Equal(identField, storage.LastRetrieved.Identifier)

	assert.Equal(1, len(storage.StoredData.Data.Commands))

	storage.Reset()

	storage.DataToReturn.Commands = append(storage.DataToReturn.Commands, storagemodels.CustomCommandsPluginCommand{
		Command:   customCommand,
		Text:      message,
		ChannelID: postToPlugin.ChannelID,
	})

	storage.DataToReturn.Commands = append(storage.DataToReturn.Commands, storagemodels.CustomCommandsPluginCommand{
		Command:   "commandNotToRemove",
		Text:      message,
		ChannelID: postToPlugin.ChannelID,
	})

	api.Reset()
	postToPlugin.User.Name = " some other user"
	postToPlugin.Content = "!customcommand remove " + customCommand
	expectedPostFromPlugin.Content = fmt.Sprintf("Custom command '%s' removed.", customCommand)

	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.User.Name = userName
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	assert.Equal(p.BotID, storage.LastRetrieved.BotID)
	assert.Equal(p.PluginID, storage.LastRetrieved.PluginID)
	assert.Equal(identField, storage.LastRetrieved.Identifier)

	assert.Equal(1, len(storage.StoredData.Data.Commands))
	assert.Equal(storage.DataToReturn.Commands[1], storage.StoredData.Data.Commands[0])

	storage.ErrorToReturn = errors.New("Some error")
	postToPlugin.Content = "!customcommand add " + customCommand + " " + message
	expectedPostFromPlugin.Content = "Could not add custom command. Please try again later."
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	postToPlugin.Content = "!customcommand remove " + customCommand
	expectedPostFromPlugin.Content = "Could not remove custom command. Please try again later."
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestCustomCommandsPlugin_UseCustomCommand(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	storage := &MockStorage{}
	api := plugin.MockAPI{Storage: storage}
	api.ProvidedFeatures = providedFeatures
	err = p.SetAPI(&api)
	assert.NoError(err)

	channel1 := "some channel"
	channel2 := "some other channel"

	ccommand1 := "someCustomCommand"
	ccommandOtherChannel := "someOtherChannelCommand"

	messageChannel1 := "some message"
	messageChannel2 := "message to other channel"
	otherMessage := "some other message"

	api.Reset()
	postToPlugin := model.Post{
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		IsPrivate: false,
	}

	storage.DataToReturn.Commands = append(storage.DataToReturn.Commands, storagemodels.CustomCommandsPluginCommand{
		Command:   ccommand1,
		Text:      messageChannel1,
		ChannelID: channel1,
	})
	storage.DataToReturn.Commands = append(storage.DataToReturn.Commands, storagemodels.CustomCommandsPluginCommand{
		Command:   ccommand1,
		Text:      messageChannel2,
		ChannelID: channel2,
	})
	storage.DataToReturn.Commands = append(storage.DataToReturn.Commands, storagemodels.CustomCommandsPluginCommand{
		Command:   ccommandOtherChannel,
		Text:      otherMessage,
		ChannelID: channel2,
	})

	expectedPostFromPlugin := model.Post{
		ChannelID: channel1,
		Content:   messageChannel1,
	}
	postToPlugin.ChannelID = channel1
	postToPlugin.Content = "!" + ccommand1
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled, "1. test failed")
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost, "1. test failed")

	api.Reset()
	postToPlugin.ChannelID = channel1
	postToPlugin.Content = "!" + ccommandOtherChannel
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled, "2. test failed")

	api.Reset()
	expectedPostFromPlugin = model.Post{
		ChannelID: channel2,
		Content:   otherMessage,
	}
	postToPlugin.ChannelID = channel2
	postToPlugin.Content = "!" + ccommandOtherChannel
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled, "3. test failed")
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost, "3. test failed")

	api.Reset()
	expectedPostFromPlugin = model.Post{
		ChannelID: channel2,
		Content:   messageChannel2,
	}
	postToPlugin.ChannelID = channel2
	postToPlugin.Content = "!" + ccommand1
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled, "4. test failed")
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost, "4. test failed")
}

func TestCustomCommandsPlugin_NoStorage(t *testing.T) {
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

	customCommand := "someCustomCommand"
	message := "some message"

	api.Reset()
	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "!customcommand add " + customCommand + " " + message,
		IsPrivate: false,
	}
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Could not add custom command. Please try again later.",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}
