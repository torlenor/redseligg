package providers

import (
	"context"

	"github.com/torlenor/redseligg/platform"
)

type MockBot struct {
	plugins []platform.BotPlugin
}

func (b *MockBot) Start() {}
func (b *MockBot) Stop()  {}

func (b *MockBot) Run(ctx context.Context) error {
	return nil
}

func (b *MockBot) AddPlugin(plugin platform.BotPlugin) {
	b.plugins = append(b.plugins, plugin)
}

func (b *MockBot) GetInfo() platform.BotInfo {
	return platform.BotInfo{}
}
