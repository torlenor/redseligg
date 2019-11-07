package config

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

// API holds the API settings for the AbyleBotter configuration
type API struct {
	Enabled bool `toml:"enabled"`
	// IP Address the REST API listens on
	// If empty or non-existing listen on all interfaces
	IP string `toml:"ip"`
	// Port the REST API listens on
	Port string `toml:"port"`
}

// General holds the general settings for the AbyleBotter configuration
type General struct {
	API API `toml:"api"`
}

// Plugins holds the plugins part of the AbyleBotter configuration
type Plugins struct {
	Echo struct {
		Enabled      bool `toml:"enabled"`
		OnlyWhispers bool `toml:"onlywhispers"`
	} `toml:"echo"`
	HTTPPing struct {
		Enabled bool `toml:"enabled"`
	} `toml:"httpping"`
	Random struct {
		Enabled bool `toml:"enabled"`
	} `toml:"random"`
}

// DiscordConfig contains config related to the Discord component
type DiscordConfig struct {
	Enabled bool    `toml:"enabled"`
	ID      string  `toml:"id"`
	Token   string  `toml:"token"`
	Secret  string  `toml:"secret"`
	Plugins Plugins `toml:"plugins"`
}

// MatrixConfig contains config related to the Matrix component
type MatrixConfig struct {
	Enabled  bool    `toml:"enabled"`
	Server   string  `toml:"server"`
	Username string  `toml:"username"`
	Password string  `toml:"password"`
	Token    string  `toml:"token"`
	Plugins  Plugins `toml:"plugins"`
}

// MattermostConfig contains config related to the Mattermost component
type MattermostConfig struct {
	Enabled  bool    `toml:"enabled"`
	Server   string  `toml:"server"`
	Username string  `toml:"username"`
	Password string  `toml:"password"`
	UseToken bool    `toml:"usetoken"`
	Token    string  `toml:"token"`
	Plugins  Plugins `toml:"plugins"`
}

// SlackConfig contains config related to the Mattermost component
type SlackConfig struct {
	Enabled   bool    `toml:"enabled" json:"enabled"`
	Workspace string  `toml:"workspace" json:"workspace"`
	Token     string  `toml:"token" json:"token"`
	Plugins   Plugins `toml:"plugins"`
}

// BotsConfig contains the complete config related to the bots
type BotsConfig struct {
	Discord    DiscordConfig    `toml:"discord"`
	Matrix     MatrixConfig     `toml:"matrix"`
	Mattermost MattermostConfig `toml:"mattermost"`
	Slack      SlackConfig      `toml:"slack"`
}

// Config holds the complete AbyleBotter config
type Config struct {
	General General    `toml:"general"`
	Bots    BotsConfig `toml:"bots"`
}

// Parse a provided reader and return the config if successful
func Parse(rd io.Reader) (Config, error) {
	var cfg Config
	if _, err := toml.DecodeReader(rd, &cfg); err != nil {
		return Config{}, fmt.Errorf("Not able to parse config: %s", err)
	}

	return cfg, nil
}

// ParseFromFile and return the config if successful
func ParseFromFile(fileName string) (Config, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return Config{}, fmt.Errorf("Not able to open file %s: %s", fileName, err)
	}

	return Parse(file)
}
