package versionplugin

import (
	"github.com/torlenor/redseligg/model"
)

// OnRun implements the hook from the Boot
func (p *VersionPlugin) OnRun() {
	p.API.RegisterCommand(p, "version")
}

// OnCommand implements the hook from the Bot
func (p *VersionPlugin) OnCommand(cmd string, content string, post model.Post) {
	versionPost := post
	versionPost.Content = p.API.GetVersion()
	_, err := p.API.CreatePost(versionPost)
	if err != nil {
		p.API.LogError("VersionPlugin: Error sending message: " + err.Error())
	}
}
