package memorybackend

import (
	"errors"
)

// ErrNotFound is returned when no data could be found
var ErrNotFound = errors.New("MemoryBackend: Could not find requested data")

// GetPluginData returns stored data of an Plugin identified by botID, pluginID and a unique identifer
// which the plugin provides. The plugin has to interpret the returned data itself.
func (b *MemoryBackend) GetPluginData(botID, pluginID, identifier string) (interface{}, error) {
	if d, ok := b.storage[botID][pluginID][identifier]; ok {
		return d, nil
	}
	return nil, ErrNotFound
}
