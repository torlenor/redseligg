package memorystorage

import "github.com/torlenor/abylebotter/storage"

// DeleteQuotesPluginQuote returns a QuotesPluginQuote.
func (b *MemoryStorage) DeleteQuotesPluginQuote(botID, pluginID, identifier string) error {
	if q, ok := b.storage[botID][pluginID]; ok {
		delete(q, identifier)
		return nil
	}
	return storage.ErrNotFound
}
