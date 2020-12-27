package rssplugin

import (
	"errors"
	"fmt"
	"testing"

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

var command = "rss"

func TestCreateRssPlugin(t *testing.T) {
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

func TestRssPlugin_HasExpectedRequiredFeatures(t *testing.T) {
	assert := assert.New(t)

	expectedRequiredFeatures := []string{
		platform.FeatureMessagePost,
	}

	p, _ := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.Equal(expectedRequiredFeatures, p.NeededFeatures)
}

func TestRssPlugin_OnRun(t *testing.T) {
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

func TestRssPlugin_HelpTextAndInvalidCommands(t *testing.T) {
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
		Content:   p.helpText(),
		IsPrivate: false,
	}
	p.OnCommand(command, "", postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	expectedPostFromPlugin.Content = p.helpTextAdd()
	p.OnCommand(command, "add", postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	expectedPostFromPlugin.Content = p.helpTextRemove()
	p.OnCommand(command, "remove", postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestRssPlugin_AddAndRemoveRssSubscription(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	storage := &MockStorage{}
	api := plugin.MockAPI{Storage: storage}
	api.ProvidedFeatures = providedFeatures
	err = p.SetAPI(&api)
	assert.NoError(err)

	link := "http://some.thing/test.xml"

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
		Content:   fmt.Sprintf("RSS subscription for link '%s' added.", link),
		IsPrivate: false,
	}
	p.OnCommand(command, "add "+link, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	assert.Equal(p.BotID, storage.LastRetrieved.BotID)
	assert.Equal(p.PluginID, storage.LastRetrieved.PluginID)

	assert.Equal(1, len(storage.StoredSubscriptions))

	storage.Reset()

	otherLink := "http://some.other.thing/test.xml"

	expectedPostFromPlugin.Content = "RSS subscription to remove does not exist."
	p.OnCommand(command, "remove "+otherLink, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	expectedPostFromPlugin.Content = "RSS subscription to remove does not exist."
	p.OnCommand(command, "remove "+link, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	storage.DataToReturn.Subscriptions = append(storage.DataToReturn.Subscriptions, storagemodels.RssPluginSubscription{
		Link:      link,
		ChannelID: postToPlugin.ChannelID,
	})

	storage.DataToReturn.Subscriptions = append(storage.DataToReturn.Subscriptions, storagemodels.RssPluginSubscription{
		Link:      "http://some.other.thing.which.should.not.be.removed/test.xml",
		ChannelID: postToPlugin.ChannelID,
	})

	expectedPostFromPlugin.Content = fmt.Sprintf("RSS subscription for link '%s' added.", "http://some.other.thing.which.should.not.be.removed/test.xml")
	p.OnCommand(command, "add "+"http://some.other.thing.which.should.not.be.removed/test.xml", postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	expectedPostFromPlugin.Content = fmt.Sprintf("RSS subscription for link '%s' removed.", link)
	p.OnCommand(command, "remove "+link, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	assert.Equal(p.BotID, storage.LastRetrieved.BotID)
	assert.Equal(p.PluginID, storage.LastRetrieved.PluginID)

	assert.Equal(1, len(storage.StoredSubscriptions))
	assert.Equal(storage.DataToReturn.Subscriptions[1].Link, storage.StoredSubscriptions[0].Data.Link)

	storage.ErrorToReturn = errors.New("Some error")
	expectedPostFromPlugin.Content = "Could not add RSS subscription. Please try again later."
	p.OnCommand(command, "add "+link, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	expectedPostFromPlugin.Content = "Could not remove RSS subscription. Please try again later."
	p.OnCommand(command, "remove "+link, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestRssPlugin_AddAndRemoveRssSubscription_OnlyMods(t *testing.T) {
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

	link := "http://some.thing/test.xml"

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
		Content:   fmt.Sprintf("RSS subscription for link '%s' added.", link),
		IsPrivate: false,
	}
	content := "add " + link
	p.OnCommand(command, content, postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.User.Name = userName
	p.OnCommand(command, content, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	assert.Equal(1, len(storage.StoredSubscriptions))

	storage.Reset()

	storage.DataToReturn.Subscriptions = append(storage.DataToReturn.Subscriptions, storagemodels.RssPluginSubscription{
		Link:      link,
		ChannelID: postToPlugin.ChannelID,
	})

	storage.DataToReturn.Subscriptions = append(storage.DataToReturn.Subscriptions, storagemodels.RssPluginSubscription{
		Link:      "http://some.other.thing.which.should.not.be.removed/test.xml",
		ChannelID: postToPlugin.ChannelID,
	})

	expectedPostFromPlugin.Content = fmt.Sprintf("RSS subscription for link '%s' added.", "http://some.other.thing.which.should.not.be.removed/test.xml")
	p.OnCommand(command, "add "+"http://some.other.thing.which.should.not.be.removed/test.xml", postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.User.Name = " some other user"
	content = "remove " + link
	expectedPostFromPlugin.Content = fmt.Sprintf("RSS subscription for link '%s' removed.", link)

	p.OnCommand(command, content, postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)

	api.Reset()
	postToPlugin.User.Name = userName
	p.OnCommand(command, content, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	assert.Equal(p.BotID, storage.LastRetrieved.BotID)
	assert.Equal(p.PluginID, storage.LastRetrieved.PluginID)

	assert.Equal(1, len(storage.StoredSubscriptions))
	assert.Equal(storage.DataToReturn.Subscriptions[1].Link, storage.StoredSubscriptions[0].Data.Link)

	storage.ErrorToReturn = errors.New("Some error")
	content = "add " + link
	expectedPostFromPlugin.Content = "Could not add RSS subscription. Please try again later."
	p.OnCommand(command, content, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	content = "remove " + link
	expectedPostFromPlugin.Content = "Could not remove RSS subscription. Please try again later."
	p.OnCommand(command, content, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestRssPlugin_NoStorage(t *testing.T) {
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

	link := "http://some.thing/test.xml"

	api.Reset()
	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "", // not used by plugin
		IsPrivate: false,
	}
	content := "add " + link
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Could not add RSS subscription. Please try again later.",
		IsPrivate: false,
	}
	p.OnCommand(command, content, postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}
