package timedmessagesplugin

import (
	"fmt"
	"strings"
	"time"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/storage"
	"github.com/torlenor/redseligg/utils"
)

func (p *TimedMessagesPlugin) checkTimedMessages(t time.Time) {
	messages, err := p.getTimedMessages()
	if err == storage.ErrNotFound {
		return
	} else if err != nil {
		p.API.LogError(fmt.Sprintf("Unable to get timed messages: %s", err))
		return
	}

	for i, m := range messages.Messages {
		if t.Sub(m.LastSent) > m.Interval {
			p.returnMessage(m.ChannelID, m.Text)
			m.LastSent = t
			messages.Messages[i] = m
		}
	}

	err = p.storeTimedMessages(messages)
	if err != nil {
		p.API.LogError(fmt.Sprintf("Unable to store timed messages after sending: %s", err))
		return
	}
}

// OnRun is called when the platform is ready
func (p *TimedMessagesPlugin) OnRun() {
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
				p.checkTimedMessages(t)
			}
		}
	}()
}

// OnStop implements the hook from the bot
func (p *TimedMessagesPlugin) OnStop() {
	if p.ticker != nil {
		p.ticker.Stop()
		p.tickerDoneChan <- true
	}
}

// OnPost implements the hook from the Bot
func (p *TimedMessagesPlugin) OnPost(post model.Post) {
	if post.IsPrivate {
		return
	}

	msg := strings.Trim(post.Content, " ")

	if !p.cfg.OnlyMods || utils.StringSliceContains(p.cfg.Mods, post.User.Name) {
		if strings.HasPrefix(msg, "!tm ") {
			p.onCommand(post)
		} else if msg == "!tm" {
			p.returnHelp(post.ChannelID)
		}
	} else {
		p.API.LogDebug("Not parsing as command, because User " + post.User.Name + " is not part of mods")
	}
}
