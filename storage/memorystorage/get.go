package memorystorage

import (
	"fmt"

	"github.com/torlenor/abylebotter/storage"
	"github.com/torlenor/abylebotter/storagemodels"
)

// GetQuotesPluginQuote returns a QuotesPluginQuote.
func (b *MemoryStorage) GetQuotesPluginQuote(botID, pluginID, identifier string) (storagemodels.QuotesPluginQuote, error) {
	if q, ok := b.storage[botID][pluginID][identifier]; ok {
		if val, ok := q.(storagemodels.QuotesPluginQuote); ok {
			return val, nil
		}
		return storagemodels.QuotesPluginQuote{}, fmt.Errorf("Stored data is not a valid QuotesPluginQuote")
	}
	return storagemodels.QuotesPluginQuote{}, storage.ErrNotFound
}

// GetQuotesPluginQuotesList returns a QuotesPluginQuotesList.
func (b *MemoryStorage) GetQuotesPluginQuotesList(botID, pluginID, identifier string) (storagemodels.QuotesPluginQuotesList, error) {
	if q, ok := b.storage[botID][pluginID][identifier]; ok {
		if val, ok := q.(storagemodels.QuotesPluginQuotesList); ok {
			return val, nil
		}
		return storagemodels.QuotesPluginQuotesList{}, fmt.Errorf("Stored data is not a valid QuotesPluginQuotesList")
	}
	return storagemodels.QuotesPluginQuotesList{}, storage.ErrNotFound
}
