package providers

import (
	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/plugin"
)

const MockPluginType = "MockPlugin"

type MockPlugin struct {
	BotID    string
	PluginID string
}

func (m *MockPlugin) SetBotPluginID(botID, pluginID string) {
	m.BotID = botID
	m.PluginID = pluginID
}

func (m *MockPlugin) SetAPI(api plugin.API) error { return nil }

// PluginType returns the plugin type
func (m *MockPlugin) PluginType() string { return MockPluginType }

func (m *MockPlugin) OnPost(model.Post)                {}
func (m *MockPlugin) OnRun()                           {}
func (m *MockPlugin) OnStop()                          {}
func (m *MockPlugin) OnReactionAdded(model.Reaction)   {}
func (m *MockPlugin) OnReactionRemoved(model.Reaction) {}
