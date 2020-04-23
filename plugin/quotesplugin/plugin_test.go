package quotesplugin

import (
	"fmt"
	"reflect"
	"testing"

	"git.abyle.org/redseligg/botorchestrator/botconfig"
	"github.com/stretchr/testify/assert"
	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/plugin"
)

func TestCreateQuotesPlugin(t *testing.T) {
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

func TestQuotesPlugin_HelpTextAndInvalidCommands(t *testing.T) {
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
	postToPlugin.Content = "!quoteadd"
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   helpText,
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!quoteadd "
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!quotehelp"
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!quoteremove"
	expectedPostFromPlugin = model.Post{
		ChannelID: "CHANNEL ID",
		Content:   helpTextRemove,
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	api.Reset()
	postToPlugin.Content = "!quoteremove something"
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestQuotesPlugin_AddQuote(t *testing.T) {
	assert := assert.New(t)

	pluginID := "SOME_PLUGIN_ID"
	quote := "some quote"

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	p.PluginID = pluginID

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

	postToPlugin.Content = "!quoteadd " + quote
	p.OnPost(postToPlugin)
	assert.Equal(false, api.WasCreatePostCalled)
	if !assert.Equal(2, len(storage.Stored)) {
		t.FailNow()
	}

	actualData := storage.Stored[0]
	assert.Equal(pluginID, actualData.PluginID)
	assert.Greater(len(actualData.Identifier), 0)
	assert.Equal(quote, actualData.Data)

	actualList := storage.Stored[1]
	assert.Equal(pluginID, actualList.PluginID)
	assert.Equal(LIST_IDENTIFIER, actualList.Identifier)
	fmt.Printf("%s", reflect.TypeOf(actualList.Data))
	if list, ok := actualList.Data.(quotesList); ok {
		assert.Equal(1, len(list))
	} else {
		assert.FailNow("Stored list not a quotesList")
	}
}
