package quotesplugin

import (
	"fmt"
	"testing"
	"time"

	"git.abyle.org/redseligg/botorchestrator/botconfig"
	"github.com/stretchr/testify/assert"
	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/plugin"
	"github.com/torlenor/abylebotter/storagemodels"
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
	storage := MockStorage{}
	p.SetAPI(&api, &storage)
}

func TestQuotesPlugin_HelpTextAndInvalidCommands(t *testing.T) {
	assert := assert.New(t)

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	api := plugin.MockAPI{}
	storage := MockStorage{}
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

	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.000Z"
		t, _ := time.Parse(layout, str)
		return t
	}

	pluginID := "SOME_PLUGIN_ID"
	quote := storagemodels.QuotesPluginQuote{
		Author:    "USER 1",
		Added:     now(),
		AuthorID:  "SOME USER ID",
		ChannelID: "CHANNEL ID",
		Text:      "some quote",
	}

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	p.PluginID = pluginID

	api := plugin.MockAPI{}
	storage := MockStorage{}
	p.SetAPI(&api, &storage)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}

	postToPlugin.Content = "!quoteadd " + quote.Text
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Successfully added quote #1",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	if !assert.Equal(1, len(storage.StoredQuotes)) {
		t.FailNow()
	}
	if !assert.Equal(1, len(storage.StoredQuotesList)) {
		t.FailNow()
	}

	actualData := storage.StoredQuotes[0]
	assert.Equal(pluginID, actualData.PluginID)
	assert.Greater(len(actualData.Identifier), 0)
	assert.Equal(quote, actualData.Data)

	actualList := storage.StoredQuotesList[0]
	assert.Equal(pluginID, actualList.PluginID)
	assert.Equal(identFieldList, actualList.Identifier)
	assert.Equal(1, len(actualList.Data.UUIDs))
}

func TestQuotesPlugin_AddQuoteFail(t *testing.T) {
	assert := assert.New(t)

	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.000Z"
		t, _ := time.Parse(layout, str)
		return t
	}

	pluginID := "SOME_PLUGIN_ID"
	quote := storagemodels.QuotesPluginQuote{
		Author:    "USER 1",
		Added:     now(),
		AuthorID:  "SOME USER ID",
		ChannelID: "CHANNEL ID",
		Text:      "some quote",
	}

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	p.PluginID = pluginID

	api := plugin.MockAPI{}
	storage := MockStorage{}
	storage.ErrorToReturn = fmt.Errorf("Some error")
	p.SetAPI(&api, &storage)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "MESSAGE CONTENT",
		IsPrivate: false,
	}

	postToPlugin.Content = "!quoteadd " + quote.Text
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Error storing quote. Try again later!",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestQuotesPlugin_GetQuote(t *testing.T) {
	assert := assert.New(t)

	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.000Z"
		t, _ := time.Parse(layout, str)
		return t
	}

	pluginID := "SOME_PLUGIN_ID"
	quote := storagemodels.QuotesPluginQuote{
		Author:    "USER 1",
		Added:     now(),
		AuthorID:  "SOME USER ID",
		ChannelID: "CHANNEL ID",
		Text:      "some quote",
	}

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	p.PluginID = pluginID

	api := plugin.MockAPI{}
	storage := MockStorage{}
	storage.QuoteDataToReturn = quote
	storage.QuotesListDataToReturn = storagemodels.QuotesPluginQuotesList{
		UUIDs: []string{"some identifier"},
	}

	p.SetAPI(&api, &storage)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "!quote",
		IsPrivate: false,
	}

	year, month, day := now().Date()
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   fmt.Sprintf(`1. "%s" - %d-%d-%d, added by %s`, quote.Text, year, month, day, postToPlugin.User.Name),
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestQuotesPlugin_GetQuote_Number(t *testing.T) {
	assert := assert.New(t)

	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.000Z"
		t, _ := time.Parse(layout, str)
		return t
	}

	pluginID := "SOME_PLUGIN_ID"

	quote2 := storagemodels.QuotesPluginQuote{
		Author:    "USER 1",
		Added:     now(),
		AuthorID:  "SOME USER ID",
		ChannelID: "CHANNEL ID",
		Text:      "some other quote",
	}

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	p.PluginID = pluginID

	api := plugin.MockAPI{}
	storage := MockStorage{}
	storage.QuoteDataToReturn = quote2
	storage.QuotesListDataToReturn = storagemodels.QuotesPluginQuotesList{
		UUIDs: []string{"some identifier", "some other identifier"},
	}

	p.SetAPI(&api, &storage)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "!quote 2",
		IsPrivate: false,
	}

	year, month, day := now().Date()
	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   fmt.Sprintf(`2. "%s" - %d-%d-%d, added by %s`, quote2.Text, year, month, day, postToPlugin.User.Name),
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)
}

func TestQuotesPlugin_RemoveQuote(t *testing.T) {
	assert := assert.New(t)

	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.000Z"
		t, _ := time.Parse(layout, str)
		return t
	}

	pluginID := "SOME_PLUGIN_ID"
	botID := "SOME BOT ID"
	quote := storagemodels.QuotesPluginQuote{
		Author:    "USER 1",
		Added:     now(),
		AuthorID:  "SOME USER ID",
		ChannelID: "CHANNEL ID",
		Text:      "some quote",
	}

	p, err := New(botconfig.PluginConfig{Type: PLUGIN_TYPE})
	assert.NoError(err)
	assert.Equal(nil, p.API)

	p.PluginID = pluginID
	p.BotID = botID

	api := plugin.MockAPI{}
	storage := MockStorage{}
	storage.QuoteDataToReturn = quote
	storage.QuotesListDataToReturn = storagemodels.QuotesPluginQuotesList{
		UUIDs: []string{"some identifier", "some other identifier"},
	}
	p.SetAPI(&api, &storage)

	postToPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Channel:   "SOME CHANNEL",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "!quoteremove 2",
		IsPrivate: false,
	}

	expectedPostFromPlugin := model.Post{
		ChannelID: "CHANNEL ID",
		Content:   "Successfully removed quote #2",
		IsPrivate: false,
	}
	p.OnPost(postToPlugin)
	assert.Equal(true, api.WasCreatePostCalled)
	assert.Equal(expectedPostFromPlugin, api.LastCreatePostPost)

	if !assert.Equal(1, len(storage.StoredQuotesList)) {
		t.FailNow()
	}

	assert.Equal(botID, storage.LastDeleted.BotID)
	assert.Equal(pluginID, storage.LastDeleted.PluginID)
	assert.Equal("some other identifier", storage.LastDeleted.Identifier)
}
