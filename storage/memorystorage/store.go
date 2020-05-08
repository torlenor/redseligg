package memorystorage

import (
	"github.com/torlenor/redseligg/storagemodels"
)

// StoreQuotesPluginQuote takes a QuotesPluginQuote and stores it.
func (b *MemoryStorage) StoreQuotesPluginQuote(botID, pluginID, identifier string, data storagemodels.QuotesPluginQuote) error {
	if _, ok := b.storage[botID]; !ok {
		b.storage[botID] = make(pluginStorage)
	}
	if _, ok := b.storage[botID][pluginID]; !ok {
		b.storage[botID][pluginID] = make(memoryStorage)
	}
	b.storage[botID][pluginID][identifier] = data

	return nil
}

// StoreQuotesPluginQuotesList takes a QuotesPluginQuotesList and stores it.
func (b *MemoryStorage) StoreQuotesPluginQuotesList(botID, pluginID, identifier string, data storagemodels.QuotesPluginQuotesList) error {
	if _, ok := b.storage[botID]; !ok {
		b.storage[botID] = make(pluginStorage)
	}
	if _, ok := b.storage[botID][pluginID]; !ok {
		b.storage[botID][pluginID] = make(memoryStorage)
	}
	b.storage[botID][pluginID][identifier] = data

	return nil
}
