package archiveplugin

import (
	"time"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/storagemodels"
)

var now = time.Now

// OnRun is called when the platform is ready
func (p *ArchivePlugin) OnRun() {
	p.storage = p.getStorage()
	if p.storage == nil {
		p.API.LogError(ErrNoValidStorage.Error())
	}
}

// OnPost implements the hook from the Bot
func (p *ArchivePlugin) OnPost(post model.Post) {
	message := storagemodels.ArchivePluginMessage{
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
	if p.storage == nil {
		p.API.LogError(ErrNoValidStorage.Error())
		return
	}
	err := p.storage.StoreArchivePluginMessage(p.BotID, p.PluginID, ident, message)
	if err != nil {
		p.API.LogError("ArchivePlugin: Error storing message: " + err.Error())
	}
}
