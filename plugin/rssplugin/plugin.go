package rssplugin

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/torlenor/redseligg/botconfig"
	"github.com/torlenor/redseligg/platform"
	"github.com/torlenor/redseligg/plugin"
	"github.com/torlenor/redseligg/storagemodels"
)

const (
	PLUGIN_TYPE    = "rss"
	PLUGIN_COMMAND = "rss"
)

// ErrNoValidStorage is set when the provided storage does not implement the correct functions
var ErrNoValidStorage = errors.New("No valid storage set")

type writer interface {
	StoreRssPluginSubscription(botID, pluginID, identifier string, data storagemodels.RssPluginSubscription) error
}

type reader interface {
	GetRssPluginSubscriptions(botID, pluginID string) (storagemodels.RssPluginSubscriptions, error)
}

type deleter interface {
	DeleteRssPluginSubscription(botID, pluginID, identifier string) error
}

type updater interface {
	UpdateRssPluginSubscription(botID, pluginID, identifier string, data storagemodels.RssPluginSubscription) error
}

type readerWriterDeleterUpdater interface {
	reader
	writer
	deleter
	updater
}

// RssPlugin is a plugin that subscribes to RSS feeds and posts them as messages.
type RssPlugin struct {
	plugin.RedseliggPlugin

	cfg config

	storage readerWriterDeleterUpdater

	ticker         *time.Ticker
	tickerDoneChan chan bool
}

// New returns a new RssPlugin
func New(pluginConfig botconfig.PluginConfig) (*RssPlugin, error) {
	cfg, err := parseConfig(pluginConfig)
	if err != nil {
		return nil, err
	}

	ep := RssPlugin{
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
func (p *RssPlugin) getStorage() readerWriterDeleterUpdater {
	if s, ok := p.API.GetStorage().(readerWriterDeleterUpdater); ok {
		return s
	}
	return nil
}

func generateIdentifier() string {
	return uuid.New().String()
}
