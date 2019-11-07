package providers

import (
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTomlBotProvider_GetBotConfig(t *testing.T) {
	assert := assert.New(t)

	expectedTomlBotConfig := tomlBotConfig{
		Bots: BotConfigs{
			"bot1": BotConfig{
				Type:   "SomePlatform",
				Config: "something",
				Plugins: PluginConfigs{
					"plugin1": PluginConfig{
						Type:   "plugintype",
						Config: "somePluginConfig",
					},
				},
			},
		},
	}

	tomlBotConfigProvider := TomlBotConfigProvider{
		cfg: expectedTomlBotConfig,
	}

	actualCfg, err := tomlBotConfigProvider.GetBotConfig("bot1")
	assert.NoError(err)
	assert.Equal(expectedTomlBotConfig.Bots["bot1"], actualCfg)

	_, err = tomlBotConfigProvider.GetBotConfig("unknownbot")
	assert.Error(err)
}

func TestTomlBotProvider_ParseTomlBotConfigFromFile(t *testing.T) {
	assert := assert.New(t)

	_, filename, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		log.Fatal(err)
	}

	expectedTomlBotConfig := tomlBotConfig{
		Bots: BotConfigs{
			"slack_dev": BotConfig{
				Type: "slack",
				Config: map[string]interface{}{
					"workspace": "something",
					"token":     "token_goes_here",
				},
				Plugins: PluginConfigs{
					"1": PluginConfig{
						Type: "echo",
						Config: map[string]interface{}{
							"onlywhispers": false,
						},
					},
					"3": PluginConfig{
						Type: "roll",
					},
				},
			},
			"mm_dev": BotConfig{
				Type: "mattermost",
				Config: map[string]interface{}{
					"server":   "https://server.com",
					"username": "username_goes_here",
					"password": "password_goes_here",
				},
				Plugins: PluginConfigs{
					"2": PluginConfig{
						Type: "httpping",
					},
					"3": PluginConfig{
						Type: "roll",
					},
				},
			},
		},
	}

	tomlBotConfigProvider, err := ParseTomlBotConfigFromFile(dir + "/../test/testdata/bots.toml")
	assert.NoError(err)
	assert.Equal(expectedTomlBotConfig, tomlBotConfigProvider.cfg)
}

func TestTomlBotProvider_ParseTomlBotConfig_Corrupt(t *testing.T) {
	assert := assert.New(t)

	_, filename, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		log.Fatal(err)
	}

	_, err = ParseTomlBotConfigFromFile(dir + "/../test/testdata/bots_corrupt.toml")
	assert.Error(err)
}

func TestTomlBotProvider_ParseTomlBotConfigFromFile_Does_Not_Exist(t *testing.T) {
	assert := assert.New(t)

	_, filename, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		log.Fatal(err)
	}

	_, err = ParseTomlBotConfigFromFile(dir + "/../test/testdata/bots_does_not_exist.toml")
	assert.Error(err)
}
