package echoplugin

import (
	"strings"

	"github.com/torlenor/redseligg/model"
)

// OnRun implements the hool from the Boot
func (p *EchoPlugin) OnRun() {
	p.API.RegisterCommand(p, "echo")
}

// OnCommand implements the hook from the Bot
func (p *EchoPlugin) OnCommand(cmd string, content string, post model.Post) {
	msg := strings.Trim(content, " ")
	if len(msg) == 0 {
		return
	}
	if !p.onlyOnWhisper || post.IsPrivate {
		echo := post
		echo.Content = msg
		p.API.CreatePost(echo)
	}
}
