package config

import (
	"fmt"
)

// PluginConfig holds the configuration for one plugin
type PluginConfig struct {
	Type   string      `toml:"type"`
	Config interface{} `toml:"config"`
}

// PluginConfigs holds a collection of PluginConfigs identified by an unique id
type PluginConfigs map[string]PluginConfig

// BotConfig holds the configuration for one bot
type BotConfig struct {
	Type    string        `toml:"type"`
	Config  interface{}   `toml:"config"`
	Plugins PluginConfigs `toml:"plugins"`
}

// BotConfigs holds a collection of BotConfigs identified by an unique id
type BotConfigs map[string]BotConfig

// AsSlackConfig converts the config to a SlackConfig
func (c *BotConfig) AsSlackConfig() (SlackConfig, error) {
	if c.Type != "slack" {
		return SlackConfig{}, fmt.Errorf("Not a slack config")
	}

	var configMap map[string]interface{}
	var ok bool
	if configMap, ok = c.Config.(map[string]interface{}); !ok {
		return SlackConfig{}, fmt.Errorf("Cannot convert config")
	}

	var workspace string
	var token string

	if workspace, ok = configMap["workspace"].(string); !ok {
		return SlackConfig{}, fmt.Errorf("Cannot convert config, missing/unconvertible workspace")
	}

	if token, ok = configMap["token"].(string); !ok {
		return SlackConfig{}, fmt.Errorf("Cannot convert config, missing/unconvertible token")
	}

	slackCfg := SlackConfig{
		Workspace: workspace,
		Token:     token,
	}

	return slackCfg, nil
}

// AsMattermostConfig converts the config to a MattermostConfig
func (c *BotConfig) AsMattermostConfig() (MattermostConfig, error) {
	if c.Type != "mattermost" {
		return MattermostConfig{}, fmt.Errorf("Not a mattermost config")
	}

	var configMap map[string]interface{}
	var ok bool
	if configMap, ok = c.Config.(map[string]interface{}); !ok {
		return MattermostConfig{}, fmt.Errorf("Cannot convert config")
	}

	var server string
	var username string
	var password string

	if server, ok = configMap["server"].(string); !ok {
		return MattermostConfig{}, fmt.Errorf("Cannot convert config, missing/unconvertible server")
	}

	if username, ok = configMap["username"].(string); !ok {
		return MattermostConfig{}, fmt.Errorf("Cannot convert config, missing/unconvertible username")
	}

	if password, ok = configMap["password"].(string); !ok {
		return MattermostConfig{}, fmt.Errorf("Cannot convert config, missing/unconvertible password")
	}

	mmCfg := MattermostConfig{
		Server:   server,
		Username: username,
		Password: password,
	}

	return mmCfg, nil
}
