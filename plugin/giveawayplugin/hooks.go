package giveawayplugin

import (
	"strings"
	"time"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/utils"
)

var command = "giveaway"

// OnRun implements the hook from the bot
func (p *GiveawayPlugin) OnRun() {
	p.API.RegisterCommand(p, command)

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

	p.API.UnRegisterCommand(command)
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

	subcommand, arguments := utils.ExtractSubCommandAndArgsString(content)

	if !p.cfg.OnlyMods || utils.StringSliceContains(p.cfg.Mods, post.User.Name) {
		if subcommand == "start" && len(arguments) > 0 {
			p.onCommandGStart(arguments, post)
			return
		} else if subcommand == "end" {
			p.onCommandGEnd(post)
			return
		} else if subcommand == "reroll" {
			p.onCommandGReroll(post)
			return
		} else if (subcommand == "start" && len(arguments) == 0) || subcommand == "help" {
			p.returnHelp(post.ChannelID)
			return
		}
	} else {
		p.API.LogDebug("Not parsing as command, because User " + post.User.Name + " is not part of mods")
	}
}
