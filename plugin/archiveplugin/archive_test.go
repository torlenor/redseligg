package archiveplugin

import (
	"fmt"
	"testing"
	"time"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/plugin"
	"github.com/torlenor/redseligg/storagemodels"

	"github.com/stretchr/testify/assert"
)

func TestArchivePlugin_OnPost(t *testing.T) {
	assert := assert.New(t)

	p := ArchivePlugin{
		RedseliggPlugin: plugin.RedseliggPlugin{
			BotID:    "BOT ID",
			PluginID: "PLUGIN ID",
		},
	}
	assert.Equal(nil, p.API)

	storage := &MockStorage{}
	api := plugin.MockAPI{Storage: storage}
	p.SetAPI(&api)
	p.OnRun()

	// Inject a new time.Now()
	now = func() time.Time {
		layout := "2006-01-02T15:04:05.000Z"
		str := "2018-12-22T13:00:00.000Z"
		t, _ := time.Parse(layout, str)
		return t
	}

	post := model.Post{
		ServerID:  "SERVER ID",
		Server:    "SERVER NAME",
		ChannelID: "CHANNEL ID",
		Channel:   "CHANNEL NAME",
		User:      model.User{ID: "SOME USER ID", Name: "USER 1"},
		Content:   "MESSAGE CONTENT",
	}

	p.OnPost(post)
	assert.Equal(false, api.WasCreatePostCalled)

	expectedStoredMessage := storagemodels.ArchivePluginMessage{
		TImestamp: now(),

		ServerID: post.ServerID,
		Server:   post.Server,

		ChannelID: post.ChannelID,
		Channel:   post.Channel,

		UserID:   post.User.ID,
		UserName: post.User.Name,

		Content: post.Content,

		IsPrivate: post.IsPrivate,
	}

	assert.Equal(1, len(storage.StoredMessages))
	assert.Equal("BOT ID", storage.StoredMessages[0].BotID)
	assert.Equal("PLUGIN ID", storage.StoredMessages[0].PluginID)
	assert.Equal(ident, storage.StoredMessages[0].Identifier)
	assert.Equal(expectedStoredMessage, storage.StoredMessages[0].Data)

	api.Reset()
	storage.ErrorToReturn = fmt.Errorf("Some error")
	p.OnPost(post)
	assert.Equal("ArchivePlugin: Error storing message: Some error", api.LastLoggedError)

	api.Reset()
	api.Storage = nil
	p.OnRun()
	assert.Equal(ErrNoValidStorage.Error(), api.LastLoggedError)

	p.OnPost(post)
	assert.Equal(ErrNoValidStorage.Error(), api.LastLoggedError)
}
