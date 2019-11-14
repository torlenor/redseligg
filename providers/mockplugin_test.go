package providers

import (
	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/plugin"
)

type MockPlugin struct{}

func (m *MockPlugin) SetAPI(api plugin.API) {}
func (m *MockPlugin) OnPost(model.Post)     {}
