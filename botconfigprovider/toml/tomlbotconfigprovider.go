package tomlbotconfigprovider

import (
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"

	"github.com/torlenor/abylebotter/botconfig"
)

type tomlBotConfig struct {
	Bots botconfig.BotConfigs `toml:"bots"`
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

	provider.setBotIDs()

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

func (c *TomlBotConfigProvider) setBotIDs() {
	for botID, botConfig := range c.cfg.Bots {
		botConfig.BotID = botID
		c.cfg.Bots[botID] = botConfig
	}
}

// GetBotConfig returns a config for the given ID if it exists
func (c *TomlBotConfigProvider) GetBotConfig(id string) (botconfig.BotConfig, error) {
	if cfg, ok := c.cfg.Bots[id]; ok {
		return cfg, nil
	}

	return botconfig.BotConfig{}, fmt.Errorf("Bot ID %s not known", id)
}

// GetAllBotConfigs returns all known bot configurations
func (c *TomlBotConfigProvider) GetAllBotConfigs() botconfig.BotConfigs {
	return c.cfg.Bots
}

// GetAllEnabledBotIDs returns only bot IDs for bots which are enabled
func (c *TomlBotConfigProvider) GetAllEnabledBotIDs() (botIDs []string) {
	for botID, cfg := range c.GetAllBotConfigs() {
		if cfg.Enabled {
			botIDs = append(botIDs, botID)
		}
	}

	return botIDs
}
