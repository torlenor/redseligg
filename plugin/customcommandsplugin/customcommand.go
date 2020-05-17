package customcommandsplugin

import (
	"fmt"
	"strings"

	"github.com/torlenor/redseligg/model"
)

func (p *CustomCommandsPlugin) onCustomCommand(post model.Post) {
	// TODO (#31): Do not fetch commands in CustomCommandPlugin every time
	customCommands, err := p.getCommands()
	if err != nil {
		p.API.LogError(fmt.Sprintf("Could not get custom commands from storage: %s", err))
		return
	}

	if len(post.Content) < 2 {
		p.API.LogWarn("'%s' cannot be a custom command: String too short")
		return
	}

	command := post.Content[1:]

	for _, c := range customCommands.Commands {
		if c.ChannelID != post.ChannelID {
			continue
		}

		if strings.HasPrefix(command, c.Command) {
			p.returnMessage(post.ChannelID, c.Text)
		}
	}
}
