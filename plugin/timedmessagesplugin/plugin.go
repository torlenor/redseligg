package timedmessagesplugin

import (
	"errors"
	"time"

	"github.com/torlenor/redseligg/botconfig"

	"github.com/torlenor/redseligg/platform"
	"github.com/torlenor/redseligg/plugin"
	"github.com/torlenor/redseligg/storagemodels"
)

const (
	PLUGIN_TYPE = "timedmessages"
)

// ErrNoValidStorage is set when the provided storage does not implement the correct functions
var ErrNoValidStorage = errors.New("No valid storage set")

// TODO (#37): TimedMessagesPlugin storage interface should only store one message at a time and not always overwrite all of them
type writer interface {
	StoreTimedMessagesPluginMessages(botID, pluginID, identifier string, data storagemodels.TimedMessagesPluginMessages) error
}

type reader interface {
	GetTimedMessagesPluginMessages(botID, pluginID, identifier string) (storagemodels.TimedMessagesPluginMessages, error)
}

type readerWriter interface {
	reader
	writer
}

// TimedMessagesPlugin is a plugin that posts messages automatically in an given interval.
type TimedMessagesPlugin struct {
	plugin.RedseliggPlugin

	cfg config

	storage readerWriter

	ticker         *time.Ticker
	tickerDoneChan chan bool
}

// New returns a new TimedMessagesPlugin
func New(pluginConfig botconfig.PluginConfig) (*TimedMessagesPlugin, error) {
	cfg, err := parseConfig(pluginConfig)
	if err != nil {
		return nil, err
	}

	ep := TimedMessagesPlugin{
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
func (p *TimedMessagesPlugin) getStorage() readerWriter {
	if s, ok := p.API.GetStorage().(readerWriter); ok {
		return s
	}
	return nil
}
