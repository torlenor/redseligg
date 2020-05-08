package customcommandsplugin

import "github.com/torlenor/redseligg/model"

const (
	helpText       = "Use '!customcommand add' to add and '!customcommand remove' to remove a custom command"
	helpTextAdd    = "Type '!customcommand add <customCommand> <message text>' to add the custom command '!customCommand'"
	helpTextRemove = "Type '!customcommand remove <customCommand> to remove the custom command '!customCommand'"
)

func (p *CustomCommandsPlugin) returnHelp(channelID string) {
	p.returnMessage(channelID, helpText)
}

func (p *CustomCommandsPlugin) returnHelpAdd(channelID string) {
	p.returnMessage(channelID, helpTextAdd)
}

func (p *CustomCommandsPlugin) returnHelpRemove(channelID string) {
	p.returnMessage(channelID, helpTextRemove)
}

func (p *CustomCommandsPlugin) returnMessage(channelID, msg string) {
	post := model.Post{
		ChannelID: channelID,
		Content:   msg,
	}
	p.API.CreatePost(post)
}
