package customcommandsplugin

import (
	"errors"
	"time"

	"github.com/torlenor/redseligg/botconfig"

	"github.com/torlenor/redseligg/platform"
	"github.com/torlenor/redseligg/plugin"
	"github.com/torlenor/redseligg/storagemodels"
)

const (
	PLUGIN_TYPE = "customcommands"
)

// ErrNoValidStorage is set when the provided storage does not implement the correct functions
var ErrNoValidStorage = errors.New("No valid storage set")

// TODO (#34): CustomCommandsPlugin storage interface should only store one command at a time and not always overwrite all of them
type writer interface {
	StoreCustomCommandsPluginCommands(botID, pluginID, identifier string, data storagemodels.CustomCommandsPluginCommands) error
}

type reader interface {
	GetCustomCommandsPluginCommands(botID, pluginID, identifier string) (storagemodels.CustomCommandsPluginCommands, error)
}

type readerWriter interface {
	reader
	writer
}

// CustomCommandsPlugin is a plugin that posts messages automatically in an given interval.
type CustomCommandsPlugin struct {
	plugin.RedseliggPlugin

	cfg config

	storage readerWriter

	ticker         *time.Ticker
	tickerDoneChan chan bool
}

// New returns a new CustomCommandsPlugin
func New(pluginConfig botconfig.PluginConfig) (*CustomCommandsPlugin, error) {
	cfg, err := parseConfig(pluginConfig)
	if err != nil {
		return nil, err
	}

	ep := CustomCommandsPlugin{
		RedseliggPlugin: plugin.RedseliggPlugin{
			NeededFeatures: []string{
				platform.FeatureMessagePost,
			},
			Type: PLUGIN_TYPE,
		},
		cfg: cfg,
	}

	return &ep, nil
}

// getStorage returns the correct storage if it supports the necessary
// functions.
func (p *CustomCommandsPlugin) getStorage() readerWriter {
	if s, ok := p.API.GetStorage().(readerWriter); ok {
		return s
	}
	return nil
}
