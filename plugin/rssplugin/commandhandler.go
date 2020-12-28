package rssplugin

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"time"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/storage"
	"github.com/torlenor/redseligg/storagemodels"
)

var now = time.Now

var errNotExist = errors.New("RSS subscription does not exist")

func (p *RssPlugin) getRssSubscriptions() (storagemodels.RssPluginSubscriptions, error) {
	s := p.getStorage()
	if s == nil {
		return storagemodels.RssPluginSubscriptions{}, ErrNoValidStorage
	}

	return s.GetRssPluginSubscriptions(p.BotID, p.PluginID)
}

func (p *RssPlugin) storeRssPluginSubscription(data storagemodels.RssPluginSubscription) error {
	s := p.getStorage()
	if s == nil {
		p.API.LogError(ErrNoValidStorage.Error())
		return ErrNoValidStorage
	}

	err := s.StoreRssPluginSubscription(p.BotID, p.PluginID, generateIdentifier(), data)
	if err != nil {
		p.API.LogError(fmt.Sprintf("Error storing RSS subscription: %s", err))
		return fmt.Errorf("Error storing RSS subscription: %s", err)
	}

	return nil
}

func (p *RssPlugin) updateRssPluginSubscription(data storagemodels.RssPluginSubscription) error {
	s := p.getStorage()
	if s == nil {
		p.API.LogError(ErrNoValidStorage.Error())
		return ErrNoValidStorage
	}

	err := s.UpdateRssPluginSubscription(p.BotID, p.PluginID, data.Identifier, data)
	if err != nil {
		p.API.LogError(fmt.Sprintf("Error updating RSS subscription: %s", err))
		return fmt.Errorf("Error updating RSS subscription: %s", err)
	}

	return nil
}

func (p *RssPlugin) addRssSubscription(channelID, link string) error {
	err := p.storeRssPluginSubscription(storagemodels.RssPluginSubscription{
		Link:              link,
		ChannelID:         channelID,
		LastPostedPubDate: time.Now(),
	})
	if err != nil {
		return fmt.Errorf("Could not add RSS subscription for link '%s' in channel '%s': %s", link, channelID, err)
	}

	p.API.LogTrace(fmt.Sprintf("Added RSS subscription for link '%s' for channel %s", link, channelID))

	return nil
}

func (p *RssPlugin) removeRssSubscription(channelID, link string) error {
	subscriptions, err := p.getRssSubscriptions()
	if err != nil && err != storage.ErrNotFound {
		return fmt.Errorf("Could not remove RSS subscription: %s", err)
	}

	s := p.getStorage()
	if s == nil {
		p.API.LogError(ErrNoValidStorage.Error())
		return ErrNoValidStorage
	}

	var wasRemoved bool
	for _, x := range subscriptions.Subscriptions {
		if x.ChannelID == channelID && x.Link == link {
			p.API.LogTrace(fmt.Sprintf("Removed RSS subscription for link '%s' for channel %s", link, channelID))
			s.DeleteRssPluginSubscription(p.BotID, p.PluginID, x.Identifier)
			wasRemoved = true
		}
	}

	if !wasRemoved {
		return errNotExist
	}

	return nil
}

func splitRssCommand(text string) (c string, link string, err error) {
	var re = regexp.MustCompile(`(?m)^+(add|remove) +(.*)$`)

	const cgCommand = 1
	const cgLink = 2

	matches := re.FindAllStringSubmatch(text, -1)

	if matches == nil || len(matches) < 1 {
		err = errors.New("Not a valid command")
		return
	}

	if len(matches[0]) > cgLink {
		c = matches[0][cgCommand]
		link = matches[0][cgLink]
	} else {
		err = errors.New("Not a valid command")
	}

	_, err = url.ParseRequestURI(link)
	if err != nil {
		err = fmt.Errorf("Not a valid url: %s", err)
	}

	return
}

func (p *RssPlugin) returnSubscriptionsList(channelID string) {
	subscriptions, err := p.getRssSubscriptions()
	if err != nil {
		p.API.LogError(fmt.Sprintf("Error getting RSS subscriptions list: %s", err))
		p.returnMessage(channelID, fmt.Sprintf("Error getting RSS subscriptions list: %s", err))
	}

	subscriptionsText := "RSS subscriptions for this channel:\n"
	cnt := 0
	lines := []string{}
	for _, s := range subscriptions.Subscriptions {
		if s.ChannelID == channelID {
			cnt++
			// subscriptionsText += fmt.Sprintf("%d. %s\n", cnt, s.Link)
			lines = append(lines, fmt.Sprintf("%d. %s", cnt, s.Link))
		}
	}

	for i, line := range lines {
		subscriptionsText += line
		if i < (len(lines) - 1) {
			subscriptionsText += "\n"
		}
	}

	p.returnMessage(channelID, subscriptionsText)
}

// onCommand handles a !rss command.
func (p *RssPlugin) onCommand(content string, post model.Post) {
	if content == "add" {
		p.returnHelpAdd(post.ChannelID)
		return
	} else if content == "remove" {
		p.returnHelpRemove(post.ChannelID)
		return
	}

	if content == "list" {
		p.returnSubscriptionsList(post.ChannelID)
		return
	}

	c, link, err := splitRssCommand(content)
	if err != nil {
		p.API.LogError(fmt.Sprintf("Error parsing command '%s': %s", content, err))
		p.returnHelp(post.ChannelID)
		return
	}

	switch c {
	case "add":
		err = p.addRssSubscription(post.ChannelID, link)
	case "remove":
		err = p.removeRssSubscription(post.ChannelID, link)
	}

	if err == errNotExist {
		p.returnMessage(post.ChannelID, "RSS subscription to remove does not exist.")
		return
	} else if err != nil {
		p.API.LogError(fmt.Sprintf("Could not %s RSS subscription: %s", c, err))
		p.returnMessage(post.ChannelID, fmt.Sprintf("Could not %s RSS subscription. Please try again later.", c))
		return
	}

	switch c {
	case "add":
		p.returnMessage(post.ChannelID, fmt.Sprintf("RSS subscription for link '%s' added.", link))
	case "remove":
		p.returnMessage(post.ChannelID, fmt.Sprintf("RSS subscription for link '%s' removed.", link))
	}
}
