package providers

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

type tomlBotConfig struct {
	Bots BotConfigs `toml:"bots"`
}

// TomlBotConfigProvider is a provider for bot configurations stored in a TOML file
type TomlBotConfigProvider struct {
	cfg tomlBotConfig
}

// ParseTomlBotConfig a provided reader and return the config if successful
func ParseTomlBotConfig(rd io.Reader) (*TomlBotConfigProvider, error) {
	var provider TomlBotConfigProvider
	var cfg tomlBotConfig
	if _, err := toml.DecodeReader(rd, &cfg); err != nil {
		return nil, fmt.Errorf("Not able to parse config: %s", err)
	}

	provider.cfg = cfg

	return &provider, nil
}

// ParseTomlBotConfigFromFile and return the config if successful
func ParseTomlBotConfigFromFile(fileName string) (*TomlBotConfigProvider, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("Not able to open file %s: %s", fileName, err)
	}

	return ParseTomlBotConfig(file)
}

// GetBotConfig returns a config for the given ID if it exists
func (c *TomlBotConfigProvider) GetBotConfig(id string) (BotConfig, error) {
	if cfg, ok := c.cfg.Bots[id]; ok {
		return cfg, nil
	}

	return BotConfig{}, fmt.Errorf("Bot ID %s not known", id)
}
