package giveawayplugin

import (
	"strings"
	"time"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/utils"
)

// OnRun implements the hook from the bot
func (p *GiveawayPlugin) OnRun() {
	p.API.RegisterCommand(p, "gstart")
	p.API.RegisterCommand(p, "gend")
	p.API.RegisterCommand(p, "greroll")
	p.API.RegisterCommand(p, "ghelp")

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
	if p.ticker != nil {
		p.ticker.Stop()
		p.tickerDoneChan <- true
	}
}

// OnPost implements the hook from the Bot
func (p *GiveawayPlugin) OnPost(post model.Post) {
	msg := strings.Trim(post.Content, " ")

	p.giveawaysMutex.Lock()
	defer p.giveawaysMutex.Unlock()

	if g, ok := p.runningGiveaways[post.ChannelID]; ok {
		if msg == g.word {
			g.addParticipant(post.User.ID, post.User.Name)
		}
	}
}

// OnCommand implements the hook from the Bot
func (p *GiveawayPlugin) OnCommand(cmd string, content string, post model.Post) {
	if post.IsPrivate {
		return
	}

	if !p.cfg.OnlyMods || utils.StringSliceContains(p.cfg.Mods, post.User.Name) {
		// TODO: Giveaway plugin should use !giveaway as base command and then start, end, reroll as sub commands
		if cmd == "gstart" && len(content) > 0 {
			p.onCommandGStart(content, post)
			return
		} else if cmd == "gend" {
			p.onCommandGEnd(post)
			return
		} else if cmd == "greroll" {
			p.onCommandGReroll(post)
			return
		} else if (cmd == "gstart" && len(content) == 0) || cmd == "ghelp" {
			p.returnHelp(post.ChannelID)
			return
		}
	} else {
		p.API.LogDebug("Not parsing as command, because User " + post.User.Name + " is not part of mods")
	}
}
