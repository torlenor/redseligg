package giveawayplugin

import (
	"strings"
	"time"

	"github.com/torlenor/abylebotter/model"
)

// OnRun implements the hook from the bot
func (p *GiveawayPlugin) OnRun() {
	p.ticker = time.NewTicker(1000 * time.Millisecond)
	p.tickerDoneChan = make(chan bool)

	go func() {
		for {
			select {
			case <-p.tickerDoneChan:
				return
			case t := <-p.ticker.C:
				for _, giveaway := range p.runningGiveaways {
					if giveaway.isFinished(t) {
						p.endGiveaway(giveaway)
					}
				}
			}
		}
	}()

	// TODO
	// ticker.Stop()
	// done <- true
}

// OnPost implements the hook from the Bot
func (p *GiveawayPlugin) OnPost(post model.Post) {
	if post.IsPrivate {
		return
	}

	msg := strings.Trim(post.Content, " ")
	if strings.HasPrefix(msg, "!gstart ") {
		p.onCommandGStart(post)
		return
	} else if strings.HasPrefix(msg, "!gend") {
		p.onCommandGEnd(post)
		return
	} else if strings.HasPrefix(msg, "!greroll") {
		p.onCommandGReroll(post)
		return
	}

	if g, ok := p.runningGiveaways[post.ChannelID]; ok {
		if msg == g.word {
			g.addParticipant(post.UserID, post.User)
		}
	}
}
