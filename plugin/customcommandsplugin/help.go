package customcommandsplugin

import (
	"fmt"

	"github.com/torlenor/redseligg/model"
)

const (
	helpText       = "Use '%scustomcommand add' to add and '%scustomcommand remove' to remove a custom command"
	helpTextAdd    = "Type '%scustomcommand add <customCommand> <message text>' to add the custom command '%scustomCommand'"
	helpTextRemove = "Type '%scustomcommand remove <customCommand> to remove the custom command '%scustomCommand'"
)

func (p *CustomCommandsPlugin) returnHelp(channelID string) {
	p.returnMessage(channelID, fmt.Sprintf(helpText, p.API.GetCallPrefix(), p.API.GetCallPrefix()))
}

func (p *CustomCommandsPlugin) returnHelpAdd(channelID string) {
	p.returnMessage(channelID, fmt.Sprintf(helpTextAdd, p.API.GetCallPrefix(), p.API.GetCallPrefix()))
}

func (p *CustomCommandsPlugin) returnHelpRemove(channelID string) {
	p.returnMessage(channelID, fmt.Sprintf(helpTextRemove, p.API.GetCallPrefix(), p.API.GetCallPrefix()))
}

func (p *CustomCommandsPlugin) returnMessage(channelID, msg string) {
	post := model.Post{
		ChannelID: channelID,
		Content:   msg,
	}
	p.API.CreatePost(post)
}
