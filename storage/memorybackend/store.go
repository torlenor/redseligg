package memorybackend

// StorePluginData stored arbitrary data of an Plugin identified by botID, pluginID and a unique identifer
// which the plugin provides.
func (b *MemoryBackend) StorePluginData(botID, pluginID, identifier string, data interface{}) error {
	if _, ok := b.storage[botID]; !ok {
		b.storage[botID] = make(pluginStorage)
	}
	if _, ok := b.storage[botID][pluginID]; !ok {
		b.storage[botID][pluginID] = make(storage)
	}
	b.storage[botID][pluginID][identifier] = data

	return nil
}
