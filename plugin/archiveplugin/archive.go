package archiveplugin

import (
	"errors"

	"github.com/torlenor/redseligg/plugin"
	"github.com/torlenor/redseligg/storagemodels"
)

// ErrNoValidStorage is set when the provided storage does not implement the correct functions
var ErrNoValidStorage = errors.New("No valid storage set")

var ident = "archive"

type writer interface {
	StoreArchivePluginMessage(botID, pluginID, identifier string, data storagemodels.ArchivePluginMessage) error
}

// ArchivePlugin is a plugin which records all messages,
type ArchivePlugin struct {
	plugin.RedseliggPlugin

	storage writer
}

// getStorage returns the correct storage if it supports the necessary
// functions.
func (p *ArchivePlugin) getStorage() writer {
	if s, ok := p.API.GetStorage().(writer); ok {
		return s
	}
	return nil
}
