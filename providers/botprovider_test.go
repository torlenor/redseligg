package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockConfigProvider struct {
}

func (m *mockConfigProvider) GetBotConfig(id string) (BotConfig, error) {
	return BotConfig{}, nil
}

func TestNewBotProvider(t *testing.T) {
	assert := assert.New(t)
	mc := &mockConfigProvider{}

	configProvider, err := NewBotProvider(mc)
	assert.NoError(err)
	assert.Same(mc, configProvider.botConfigs)
}

func TestBotProvider_createPlatformPlugins(t *testing.T) {
	// TODO
}

func TestBotProvider_GetBot(t *testing.T) {
	// TODO
}
