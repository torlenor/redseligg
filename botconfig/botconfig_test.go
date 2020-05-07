package botconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBotConfig_AsDiscordConfig(t *testing.T) {
	assert := assert.New(t)

	botConfig := BotConfig{
		Type: "discord",
		Config: map[string]interface{}{
			"id":     "some_id",
			"token":  "username_goes_here",
			"secret": "sectet_goes_here",
		},
	}

	expectedDConfig := DiscordConfig{
		ID:     "some_id",
		Token:  "username_goes_here",
		Secret: "sectet_goes_here",
	}

	actualDConfig, err := botConfig.AsDiscordConfig()
	assert.NoError(err)
	assert.Equal(expectedDConfig, actualDConfig)

	_, err = botConfig.AsSlackConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type:   "something_else",
		Config: map[string]interface{}{},
	}
	_, err = botConfig.AsDiscordConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "discord",
		Config: map[string]interface{}{
			"token":  "username_goes_here",
			"secret": "sectet_goes_here",
		},
	}
	_, err = botConfig.AsDiscordConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "discord",
		Config: map[string]interface{}{
			"id":     "some_id",
			"secret": "sectet_goes_here",
		},
	}
	_, err = botConfig.AsDiscordConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "discord",
		Config: map[string]interface{}{
			"id":    "some_id",
			"token": "username_goes_here",
		},
	}
	_, err = botConfig.AsDiscordConfig()
	assert.Error(err)
}

func TestBotConfig_AsMatrixConfig(t *testing.T) {
	assert := assert.New(t)

	botConfig := BotConfig{
		Type: "matrix",
		Config: map[string]interface{}{
			"server":   "https://server.com",
			"username": "username_goes_here",
			"password": "password_goes_here",
		},
	}

	expectedMConfig := MatrixConfig{
		Server:   "https://server.com",
		Username: "username_goes_here",
		Password: "password_goes_here",
	}

	actualMConfig, err := botConfig.AsMatrixConfig()
	assert.NoError(err)
	assert.Equal(expectedMConfig, actualMConfig)

	_, err = botConfig.AsSlackConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "matrix",
		Config: map[string]interface{}{
			"username": "username_goes_here",
			"password": "password_goes_here",
		},
	}
	_, err = botConfig.AsMatrixConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "matrix",
		Config: map[string]interface{}{
			"server":   "https://server.com",
			"password": "password_goes_here",
		},
	}
	_, err = botConfig.AsMatrixConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "matrix",
		Config: map[string]interface{}{
			"server":   "https://server.com",
			"username": "username_goes_here",
		},
	}
	_, err = botConfig.AsMatrixConfig()
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

	expectedMMConfig := MattermostConfig{
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

func TestBotConfig_AsSlackConfig(t *testing.T) {
	assert := assert.New(t)

	botConfig := BotConfig{
		Type: "slack",
		Config: map[string]interface{}{
			"workspace": "something",
			"token":     "token_goes_here",
		},
	}

	expectedSlackConfig := SlackConfig{
		Workspace: "something",
		Token:     "token_goes_here",
	}

	actualSlackConfig, err := botConfig.AsSlackConfig()
	assert.NoError(err)
	assert.Equal(expectedSlackConfig, actualSlackConfig)

	_, err = botConfig.AsMatrixConfig()
	assert.Error(err)

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

func TestBotConfig_AsTwitchConfig(t *testing.T) {
	assert := assert.New(t)

	botConfig := BotConfig{
		Type: "twitch",
		Config: map[string]interface{}{
			"username": "some_username",
			"token":    "some_token",
			"channels": []interface{}{"channel1", "channel2"},
		},
	}

	expectedTConfig := TwitchConfig{
		Username: "some_username",
		Token:    "some_token",
		Channels: []string{"channel1", "channel2"},
	}

	actualTConfig, err := botConfig.AsTwitchConfig()
	assert.NoError(err)
	assert.Equal(expectedTConfig, actualTConfig)

	_, err = botConfig.AsSlackConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "matrix",
		Config: map[string]interface{}{
			"username": "username_goes_here",
			"password": "password_goes_here",
		},
	}
	_, err = botConfig.AsMatrixConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "matrix",
		Config: map[string]interface{}{
			"server":   "https://server.com",
			"password": "password_goes_here",
		},
	}
	_, err = botConfig.AsMatrixConfig()
	assert.Error(err)

	botConfig = BotConfig{
		Type: "matrix",
		Config: map[string]interface{}{
			"server":   "https://server.com",
			"username": "username_goes_here",
		},
	}
	_, err = botConfig.AsMatrixConfig()
	assert.Error(err)
}
