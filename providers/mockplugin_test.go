package providers

import (
	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/plugin"
)

type MockPlugin struct{}

func (m *MockPlugin) SetAPI(api plugin.API, storageAPI plugin.StorageAPI) {}

func (m *MockPlugin) OnPost(model.Post)                {}
func (m *MockPlugin) OnRun()                           {}
func (m *MockPlugin) OnStop()                          {}
func (m *MockPlugin) OnReactionAdded(model.Reaction)   {}
func (m *MockPlugin) OnReactionRemoved(model.Reaction) {}
