package giveawayplugin

import (
	"strings"
	"time"

	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/utils"
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
}

// OnStop implements the hook from the bot
func (p *GiveawayPlugin) OnStop() {
	p.ticker.Stop()
	if p.ticker != nil {
		p.ticker.Stop()
		p.tickerDoneChan <- true
	}
}

// OnPost implements the hook from the Bot
func (p *GiveawayPlugin) OnPost(post model.Post) {
	if post.IsPrivate {
		return
	}

	msg := strings.Trim(post.Content, " ")
	if !p.cfg.OnlyMods || utils.StringSliceContains(p.cfg.Mods, post.User.Name) {
		if strings.HasPrefix(msg, "!gstart ") {
			p.onCommandGStart(post)
			return
		} else if msg == "!gend" {
			p.onCommandGEnd(post)
			return
		} else if msg == "!greroll" {
			p.onCommandGReroll(post)
			return
		} else if msg == "!gstart" || msg == "!ghelp" {
			p.returnHelp(post.ChannelID)
			return
		}
	} else {
		p.API.LogDebug("Not parsing as command, because User " + post.User.Name + " is not part of mods")
	}

	p.giveawaysMutex.Lock()
	defer p.giveawaysMutex.Unlock()

	if g, ok := p.runningGiveaways[post.ChannelID]; ok {
		if msg == g.word {
			g.addParticipant(post.User.ID, post.User.Name)
		}
	}
}
