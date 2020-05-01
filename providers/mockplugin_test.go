package providers

import (
	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/plugin"
)

type MockPlugin struct {
	BotID    string
	PluginID string
}

func (m *MockPlugin) SetBotPluginID(botID, pluginID string) {
	m.BotID = botID
	m.PluginID = pluginID
}

func (m *MockPlugin) SetAPI(api plugin.API) {}

func (m *MockPlugin) OnPost(model.Post)                {}
func (m *MockPlugin) OnRun()                           {}
func (m *MockPlugin) OnStop()                          {}
func (m *MockPlugin) OnReactionAdded(model.Reaction)   {}
func (m *MockPlugin) OnReactionRemoved(model.Reaction) {}
