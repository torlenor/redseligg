package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/torlenor/abylebotter/config"
)

func TestNewBotProvider(t *testing.T) {
	assert := assert.New(t)
	mc := &mockConfigProvider{}
	mbf := &MockBotFactory{}
	mpf := &MockPluginFactory{}

	botProvider, err := NewBotProvider(mc, mbf, mpf)
	assert.NoError(err)
	assert.Same(mc, botProvider.botConfigs)
	assert.Same(mbf, botProvider.botFactory)
	assert.Same(mpf, botProvider.pluginFactory)
}

func TestBotProvider_GetBot(t *testing.T) {
	assert := assert.New(t)
	mc := &mockConfigProvider{}
	mbf := &MockBotFactory{
		bot: MockBot{},
	}
	mpf := &MockPluginFactory{}

	botProvider, err := NewBotProvider(mc, mbf, mpf)
	assert.NoError(err)

	bot, err := botProvider.GetBot("mockSlackID")
	assert.NoError(err)
	if val, ok := bot.(*MockBot); ok {
		assert.Same(val, &mbf.bot)
	} else {
		assert.Fail("Did not get a MockBot")
	}

	bot, err = botProvider.GetBot("mockMattermostID")
	assert.NoError(err)
	if val, ok := bot.(*MockBot); ok {
		assert.Same(val, &mbf.bot)
	} else {
		assert.Fail("Did not get a MockBot")
	}

	bot, err = botProvider.GetBot("Unknown")
	assert.Error(err)
	assert.Nil(bot)

	bot, err = botProvider.GetBot("mockSomeOtherPlatformID")
	assert.Error(err)
	assert.Nil(bot)

}

func TestBotProvider_GetBot_WithPlugins(t *testing.T) {
	plConfig := config.PluginConfigs{
		"1": config.PluginConfig{
			Type: "mockEcho",
			Config: map[string]interface{}{
				"onlywhispers": false,
			},
		},
		"2": config.PluginConfig{
			Type: "somethingWhichFails",
		},
		"3": config.PluginConfig{
			Type: "mockRoll",
		},
	}

	assert := assert.New(t)
	mc := &mockConfigProvider{
		pluginsConfig: plConfig,
	}
	mbf := &MockBotFactory{
		bot: MockBot{},
	}
	mpf := &MockPluginFactory{}

	botProvider, err := NewBotProvider(mc, mbf, mpf)
	assert.NoError(err)

	bot, err := botProvider.GetBot("mockSlackID")
	assert.NoError(err)
	if val, ok := bot.(*MockBot); ok {
		assert.Same(val, &mbf.bot)
	} else {
		assert.Fail("Did not get a MockBot")
	}

	assert.Equal(2, len(mbf.bot.plugins))
	for key := range mbf.bot.plugins {
		assert.Same(&mpf.plugin, mbf.bot.plugins[key].(*MockPlugin))
	}
}
