package rssplugin

import (
	"fmt"

	"github.com/torlenor/redseligg/model"
)

const (
	helpText       = "Use `%srss add` to add or `%srss remove` to remove an RSS subscription for this channel. Use `%srss list` to list all the subscriptions for this channel."
	helpTextAdd    = "Type `%srss add <link_to_rss_feed>` to add an RSS subscription for the current channel"
	helpTextRemove = "Type `%srss remove <link_to_rss_feed>` to remove an RSS subscription for the current channel"
)

func (p *RssPlugin) helpText() string {
	return fmt.Sprintf(helpText, p.API.GetCallPrefix(), p.API.GetCallPrefix(), p.API.GetCallPrefix())
}

func (p *RssPlugin) helpTextAdd() string {
	return fmt.Sprintf(helpTextAdd, p.API.GetCallPrefix())
}

func (p *RssPlugin) helpTextRemove() string {
	return fmt.Sprintf(helpTextRemove, p.API.GetCallPrefix())
}

func (p *RssPlugin) returnHelp(channelID string) {
	p.returnMessage(channelID, p.helpText())
}

func (p *RssPlugin) returnHelpAdd(channelID string) {
	p.returnMessage(channelID, p.helpTextAdd())
}

func (p *RssPlugin) returnHelpRemove(channelID string) {
	p.returnMessage(channelID, p.helpTextRemove())
}

func (p *RssPlugin) returnMessage(channelID, msg string) {
	post := model.Post{
		ChannelID: channelID,
		Content:   msg,
	}
	p.API.CreatePost(post)
}
