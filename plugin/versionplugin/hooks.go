package versionplugin

import (
	"fmt"
	"strings"

	"github.com/torlenor/redseligg/model"
)

// OnPost implements the hook from the Bot
func (p *VersionPlugin) OnPost(post model.Post) {
	msg := strings.Trim(post.Content, " ")
	if strings.HasPrefix(msg, "!version") {
		versionPost := post
		versionPost.Content = p.API.GetVersion()
		p.API.LogTrace(fmt.Sprintf("Echoing version back to Channel = %s, content = %s", versionPost.Channel, versionPost.Content))
		_, err := p.API.CreatePost(versionPost)
		if err != nil {
			p.API.LogError("VersionPlugin: Error sending message: " + err.Error())
		}
	}
}
