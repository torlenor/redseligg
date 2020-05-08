package tomlbotconfigprovider

import (
	"log"
	"path/filepath"
	"reflect"
	"runtime"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/torlenor/redseligg/botconfig"
)

func TestTomlBotProvider_GetBotConfig(t *testing.T) {
	assert := assert.New(t)

	expectedTomlBotConfig := tomlBotConfig{
		Bots: botconfig.BotConfigs{
			"bot1": botconfig.BotConfig{
				Type:   "SomePlatform",
				Config: map[string]interface{}{"workspace": "something"},
				Plugins: botconfig.PluginConfigs{
					"plugin1": botconfig.PluginConfig{
						Type: "plugintype",
						Config: map[string]interface{}{
							"onlywhispers": false,
						},
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
		Bots: botconfig.BotConfigs{
			"slack_dev": botconfig.BotConfig{
				BotID: "slack_dev",
				Type:  "slack",
				Config: map[string]interface{}{
					"workspace": "something",
					"token":     "token_goes_here",
				},
				Plugins: botconfig.PluginConfigs{
					"1": botconfig.PluginConfig{
						Type: "echo",
						Config: map[string]interface{}{
							"onlywhispers": false,
						},
					},
					"3": botconfig.PluginConfig{
						Type: "roll",
					},
				},
			},
			"mm_dev": botconfig.BotConfig{
				BotID: "mm_dev",
				Type:  "mattermost",
				Config: map[string]interface{}{
					"server":   "https://server.com",
					"username": "username_goes_here",
					"password": "password_goes_here",
				},
				Plugins: botconfig.PluginConfigs{
					"2": botconfig.PluginConfig{
						Type: "httpping",
					},
					"3": botconfig.PluginConfig{
						Type: "roll",
					},
				},
			},
		},
	}

	tomlBotConfigProvider, err := ParseTomlBotConfigFromFile(dir + "/../../test/testdata/bots.toml")
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

	_, err = ParseTomlBotConfigFromFile(dir + "/../../test/testdata/bots_corrupt.toml")
	assert.Error(err)
}

func TestTomlBotProvider_ParseTomlBotConfigFromFile_Does_Not_Exist(t *testing.T) {
	assert := assert.New(t)

	_, filename, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		log.Fatal(err)
	}

	_, err = ParseTomlBotConfigFromFile(dir + "/../../test/testdata/bots_does_not_exist.toml")
	assert.Error(err)
}

func TestTomlBotConfigProvider_GetAllBotConfigs(t *testing.T) {
	expectedBotConfigs := tomlBotConfig{
		Bots: botconfig.BotConfigs{
			"1": botconfig.BotConfig{
				BotID: "test",
			},
			"2": botconfig.BotConfig{
				BotID: "test",
			},
		},
	}

	type fields struct {
		cfg tomlBotConfig
	}
	tests := []struct {
		name   string
		fields fields
		want   botconfig.BotConfigs
	}{
		{
			name: "Get all stored bot configs",
			fields: fields{
				cfg: expectedBotConfigs,
			},
			want: expectedBotConfigs.Bots,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &TomlBotConfigProvider{
				cfg: tt.fields.cfg,
			}
			if got := c.GetAllBotConfigs(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("TomlBotConfigProvider.GetAllBotConfigs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTomlBotConfigProvider_GetAllEnabledBotIDs(t *testing.T) {
	expectedBotConfigs := tomlBotConfig{
		Bots: botconfig.BotConfigs{
			"1": botconfig.BotConfig{
				BotID:   "test",
				Enabled: true,
			},
			"2": botconfig.BotConfig{
				BotID: "test",
			},
			"3": botconfig.BotConfig{
				BotID:   "test",
				Enabled: true,
			},
		},
	}

	type fields struct {
		cfg tomlBotConfig
	}
	tests := []struct {
		name       string
		fields     fields
		wantBotIDs []string
	}{
		{
			name: "Get all enabled bot ids",
			fields: fields{
				cfg: expectedBotConfigs,
			},
			wantBotIDs: []string{
				"1", "3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &TomlBotConfigProvider{
				cfg: tt.fields.cfg,
			}
			gotBotIDs := c.GetAllEnabledBotIDs()
			wantBotIDs := tt.wantBotIDs
			sort.Strings(gotBotIDs)
			sort.Strings(wantBotIDs)
			if !reflect.DeepEqual(gotBotIDs, wantBotIDs) {
				t.Errorf("TomlBotConfigProvider.GetAllEnabledBotIDs() = %v, want %v", gotBotIDs, wantBotIDs)
			}
		})
	}
}
