package rssplugin

import (
	"fmt"
	"time"

	"github.com/mmcdole/gofeed"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/storage"
	"github.com/torlenor/redseligg/utils"
)

func parseRssFeedFromURL(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	return feed, err
}

func (p *RssPlugin) checkRssSubscriptions(t time.Time) {
	subscriptions, err := p.getRssSubscriptions()
	if err == storage.ErrNotFound {
		return
	} else if err != nil {
		p.API.LogError(fmt.Sprintf("Unable to get rss subscriptions: %s", err))
		return
	}

	for _, m := range subscriptions.Subscriptions {
		newLastPostedPubDate := time.Now()
		wasPosted := false
		feed, err := parseRssFeedFromURL(m.Link)
		if err != nil {
			p.API.LogError(fmt.Sprintf("Unable to fetch RSS subscription for '%s' in channel %s: %s", m.Link, m.ChannelID, err))
			continue
		}
		for _, item := range feed.Items {
			if item.PublishedParsed.Sub(m.LastPostedPubDate) > 0 {
				p.API.CreatePost(model.Post{
					ChannelID: m.ChannelID,
					// TODO (#61): Based on additional optional arguments when subscribing more than just the title of the RSS feed item should be posted
					Content: feed.Title + ": " + item.Title + "\n" + "<" + item.Link + ">",
				})
				wasPosted = true
			}
		}
		if wasPosted {
			m.LastPostedPubDate = newLastPostedPubDate
			p.updateRssPluginSubscription(m)
		}
	}
}

// OnRun is called when the platform is ready
func (p *RssPlugin) OnRun() {
	p.API.RegisterCommand(p, PLUGIN_COMMAND)

	p.storage = p.getStorage()
	if p.storage == nil {
		p.API.LogError(ErrNoValidStorage.Error())
		return
	}

	p.ticker = time.NewTicker(5000 * time.Millisecond)
	p.tickerDoneChan = make(chan bool)

	go func() {
		for {
			select {
			case <-p.tickerDoneChan:
				return
			case t := <-p.ticker.C:
				p.checkRssSubscriptions(t)
			}
		}
	}()
}

// OnStop implements the hook from the bot
func (p *RssPlugin) OnStop() {
	if p.ticker != nil {
		p.ticker.Stop()
		p.tickerDoneChan <- true
	}
}

// OnCommand implements the hook from the Bot
func (p *RssPlugin) OnCommand(cmd string, content string, post model.Post) {
	if post.IsPrivate {
		return
	}

	if !p.cfg.OnlyMods || utils.StringSliceContains(p.cfg.Mods, post.User.Name) {
		if len(content) > 0 {
			p.onCommand(content, post)
		} else {
			p.returnHelp(post.ChannelID)
		}
	} else {
		p.API.LogDebug("Not parsing as command, because User " + post.User.Name + " is not part of mods")
	}
}
