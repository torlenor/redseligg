package timedmessagesplugin

import (
	"fmt"

	"github.com/torlenor/redseligg/model"
)

const (
	helpText       = "Use `%stm add` to add or `%stm remove` to remove timed messages"
	helpTextAdd    = "Type `%stm add <interval> <message text>` to add a timed message"
	helpTextRemove = "Type `%stm remove <interval> <message text>` or `%stm remove all <message text>` to remove all messages matching the <message text>"
)

func (p *TimedMessagesPlugin) helpText() string {
	return fmt.Sprintf(helpText, p.API.GetCallPrefix(), p.API.GetCallPrefix())
}

func (p *TimedMessagesPlugin) helpTextAdd() string {
	return fmt.Sprintf(helpTextAdd, p.API.GetCallPrefix())
}

func (p *TimedMessagesPlugin) helpTextRemove() string {
	return fmt.Sprintf(helpTextRemove, p.API.GetCallPrefix(), p.API.GetCallPrefix())
}

func (p *TimedMessagesPlugin) returnHelp(channelID string) {
	p.returnMessage(channelID, p.helpText())
}

func (p *TimedMessagesPlugin) returnHelpAdd(channelID string) {
	p.returnMessage(channelID, p.helpTextAdd())
}

func (p *TimedMessagesPlugin) returnHelpRemove(channelID string) {
	p.returnMessage(channelID, p.helpTextRemove())
}

func (p *TimedMessagesPlugin) returnMessage(channelID, msg string) {
	post := model.Post{
		ChannelID: channelID,
		Content:   msg,
	}
	p.API.CreatePost(post)
}
