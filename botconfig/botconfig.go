package botconfig

import (
	"fmt"
	"reflect"
)

// StorageConfig holds the configuration for a storage
type StorageConfig struct {
	Type string `toml:"type" bson:"type"`

	Config map[string]interface{} `toml:"config" bson:"config"`
}

// PluginConfig holds the configuration for one plugin
type PluginConfig struct {
	Type string `toml:"type" bson:"type"`

	Config map[string]interface{} `toml:"config" bson:"config"`
}

// PluginConfigs holds a collection of PluginConfigs identified by an unique id
type PluginConfigs map[string]PluginConfig

// BotConfig holds the configuration for one bot
type BotConfig struct {
	BotID   string `toml:"id" bson:"id"`
	Type    string `toml:"type" bson:"type"`
	Enabled bool   `toml:"enabled" bson:"enabled"`

	StorageConfig StorageConfig `toml:"storage" bson:"storage"`

	Config  map[string]interface{} `toml:"config" bson:"config"`
	Plugins PluginConfigs          `toml:"plugins" bson:"plugins"`
}

// BotConfigs holds a collection of BotConfigs identified by an unique id
type BotConfigs map[string]BotConfig

// AsDiscordConfig converts the config to a DiscordConfig
func (c *BotConfig) AsDiscordConfig() (DiscordConfig, error) {
	if c.Type != "discord" {
		return DiscordConfig{}, fmt.Errorf("Not a Discord config")
	}

	var id string
	var token string
	var secret string

	var ok bool
	if id, ok = c.Config["id"].(string); !ok {
		return DiscordConfig{}, fmt.Errorf("Cannot convert to Discord config, missing/unconvertible id")
	}

	if token, ok = c.Config["token"].(string); !ok {
		return DiscordConfig{}, fmt.Errorf("Cannot convert to Discord config, missing/unconvertible token")
	}

	if secret, ok = c.Config["secret"].(string); !ok {
		return DiscordConfig{}, fmt.Errorf("Cannot convert to Discord config, missing/unconvertible secret")
	}

	discordCfg := DiscordConfig{
		ID:     id,
		Token:  token,
		Secret: secret,
	}

	return discordCfg, nil
}

// AsMatrixConfig converts the config to a MatrixConfig
func (c *BotConfig) AsMatrixConfig() (MatrixConfig, error) {
	if c.Type != "matrix" {
		return MatrixConfig{}, fmt.Errorf("Not a Matrix config")
	}

	var server string
	var username string
	var password string

	var ok bool
	if server, ok = c.Config["server"].(string); !ok {
		return MatrixConfig{}, fmt.Errorf("Cannot convert to Matrix config, missing/unconvertible server")
	}

	if username, ok = c.Config["username"].(string); !ok {
		return MatrixConfig{}, fmt.Errorf("Cannot convert to Matrix config, missing/unconvertible username")
	}

	if password, ok = c.Config["password"].(string); !ok {
		return MatrixConfig{}, fmt.Errorf("Cannot convert to Matrix config, missing/unconvertible password")
	}

	mCfg := MatrixConfig{
		Server:   server,
		Username: username,
		Password: password,
	}

	return mCfg, nil
}

// AsMattermostConfig converts the config to a MattermostConfig
func (c *BotConfig) AsMattermostConfig() (MattermostConfig, error) {
	if c.Type != "mattermost" {
		return MattermostConfig{}, fmt.Errorf("Not a Mattermost config")
	}

	var server string
	var username string
	var password string

	var ok bool
	if server, ok = c.Config["server"].(string); !ok {
		return MattermostConfig{}, fmt.Errorf("Cannot convert to MatterMost config, missing/unconvertible server")
	}

	if username, ok = c.Config["username"].(string); !ok {
		return MattermostConfig{}, fmt.Errorf("Cannot convert to MatterMost config, missing/unconvertible username")
	}

	if password, ok = c.Config["password"].(string); !ok {
		return MattermostConfig{}, fmt.Errorf("Cannot convert to MatterMost config, missing/unconvertible password")
	}

	mmCfg := MattermostConfig{
		Server:   server,
		Username: username,
		Password: password,
	}

	return mmCfg, nil
}

// AsSlackConfig converts the config to a SlackConfig
func (c *BotConfig) AsSlackConfig() (SlackConfig, error) {
	if c.Type != "slack" {
		return SlackConfig{}, fmt.Errorf("Not a Slack config")
	}

	var workspace string
	var token string

	var ok bool
	if workspace, ok = c.Config["workspace"].(string); !ok {
		return SlackConfig{}, fmt.Errorf("Cannot convert to Slack config, missing/unconvertible workspace")
	}

	if token, ok = c.Config["token"].(string); !ok {
		return SlackConfig{}, fmt.Errorf("Cannot convert to Slack config, missing/unconvertible token")
	}

	slackCfg := SlackConfig{
		Workspace: workspace,
		Token:     token,
	}

	return slackCfg, nil
}

// AsTwitchConfig converts the config to a TwitchConfig
func (c *BotConfig) AsTwitchConfig() (TwitchConfig, error) {
	if c.Type != "twitch" {
		return TwitchConfig{}, fmt.Errorf("Not a Twitch config")
	}

	var username string
	var token string
	var channels []string

	var ok bool
	if username, ok = c.Config["username"].(string); !ok {
		return TwitchConfig{}, fmt.Errorf("Cannot convert to Twitch config, missing/unconvertible username")
	}

	if token, ok = c.Config["token"].(string); !ok {
		return TwitchConfig{}, fmt.Errorf("Cannot convert to Twitch config, missing/unconvertible token")
	}

	fmt.Printf("Type: %s\n", reflect.TypeOf(c.Config["channels"]))
	if channelsEntries, ok := c.Config["channels"].([]interface{}); ok {
		for _, channelEntry := range channelsEntries {
			if channel, ok := channelEntry.(string); ok {
				channels = append(channels, channel)
			}
		}
	} else {
		return TwitchConfig{}, fmt.Errorf("Cannot convert to Twitch config, missing/unconvertible channels")
	}

	cfg := TwitchConfig{
		Username: username,
		Token:    token,
		Channels: channels,
	}

	return cfg, nil
}
