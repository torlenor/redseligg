package httppingplugin

import (
	"fmt"

	"github.com/torlenor/redseligg/model"
)

// OnRun implements the hook from the Boot
func (p *HTTPPingPlugin) OnRun() {
	p.API.RegisterCommand(p, "httpping")
}

// OnCommand implements the hook from the Bot
func (p *HTTPPingPlugin) OnCommand(cmd string, content string, post model.Post) {
	timeMs, err := httpPing(content)

	var response string
	if err != nil {
		response = fmt.Sprintf("FAIL (%s).", err)
	} else {
		response = fmt.Sprintf("SUCCESS. Request took %d ms", timeMs)
	}

	pingReply := post
	pingReply.Content = response
	p.API.CreatePost(pingReply)
}
