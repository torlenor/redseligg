package providers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/torlenor/abylebotter/config"
)

func TestBotConfig_AsSlackConfig(t *testing.T) {
	assert := assert.New(t)

	botConfig := BotConfig{
		Type: "slack",
		Config: map[string]interface{}{
			"workspace": "something",
			"token":     "token_goes_here",
		},
	}

	expectedSlackConfig := config.SlackConfig{
		Workspace: "something",
		Token:     "token_goes_here",
	}

	actualSlackConfig, err := botConfig.AsSlackConfig()
	assert.NoError(err)
	assert.Equal(expectedSlackConfig, actualSlackConfig)

	_, err = botConfig.AsMattermostConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "slack",
		Config: map[string]interface{}{
			"workspace": "something",
		},
	}
	_, err = botConfig.AsSlackConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "slack",
		Config: map[string]interface{}{
			"token": "token_goes_here",
		},
	}
	_, err = botConfig.AsSlackConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "slack",
		Config: map[string]interface{}{
			"token": false,
		},
	}
	_, err = botConfig.AsSlackConfig()
	assert.Error(err)
}

func TestBotConfig_AsMattermostConfig(t *testing.T) {
	assert := assert.New(t)

	botConfig := BotConfig{
		Type: "mattermost",
		Config: map[string]interface{}{
			"server":   "https://server.com",
			"username": "username_goes_here",
			"password": "password_goes_here",
		},
	}

	expectedMMConfig := config.MattermostConfig{
		Server:   "https://server.com",
		Username: "username_goes_here",
		Password: "password_goes_here",
	}

	actualMMConfig, err := botConfig.AsMattermostConfig()
	assert.NoError(err)
	assert.Equal(expectedMMConfig, actualMMConfig)

	_, err = botConfig.AsSlackConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "mattermost",
		Config: map[string]interface{}{
			"username": "username_goes_here",
			"password": "password_goes_here",
		},
	}
	_, err = botConfig.AsMattermostConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "mattermost",
		Config: map[string]interface{}{
			"server":   "https://server.com",
			"password": "password_goes_here",
		},
	}
	_, err = botConfig.AsMattermostConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "mattermost",
		Config: map[string]interface{}{
			"server":   "https://server.com",
			"username": "username_goes_here",
		},
	}
	_, err = botConfig.AsMattermostConfig()
	assert.Error(err)
}
