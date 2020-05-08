package timedmessagesplugin

import "github.com/torlenor/redseligg/model"

const (
	helpText       = "Use !tm add to add and !tm remove to remove timed messages"
	helpTextAdd    = "Type !tm add <interval> <message text> to add a timed message"
	helpTextRemove = "Type !tm remove <interval> <message text> or !tm remove all <message text> to remove all messages matching the <message text>"
)

func (p *TimedMessagesPlugin) returnHelp(channelID string) {
	p.returnMessage(channelID, helpText)
}

func (p *TimedMessagesPlugin) returnHelpAdd(channelID string) {
	p.returnMessage(channelID, helpTextAdd)
}

func (p *TimedMessagesPlugin) returnHelpRemove(channelID string) {
	p.returnMessage(channelID, helpTextRemove)
}

func (p *TimedMessagesPlugin) returnMessage(channelID, msg string) {
	post := model.Post{
		ChannelID: channelID,
		Content:   msg,
	}
	p.API.CreatePost(post)
}
