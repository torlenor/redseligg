package memorystorage

import (
	"errors"
	"fmt"

	"github.com/torlenor/abylebotter/storagemodels"
)

// ErrNotFound is returned when no data could be found
var ErrNotFound = errors.New("MemoryStorage: Could not find requested data")

// GetQuotesPluginQuote returns a QuotesPluginQuote.
func (b *MemoryStorage) GetQuotesPluginQuote(botID, pluginID, identifier string) (storagemodels.QuotesPluginQuote, error) {
	if q, ok := b.storage[botID][pluginID][identifier]; ok {
		if val, ok := q.(storagemodels.QuotesPluginQuote); ok {
			return val, nil
		}
		return storagemodels.QuotesPluginQuote{}, fmt.Errorf("Stored data is not a valid QuotesPluginQuote")
	}
	return storagemodels.QuotesPluginQuote{}, ErrNotFound
}

// GetQuotesPluginQuotesList returns a QuotesPluginQuotesList.
func (b *MemoryStorage) GetQuotesPluginQuotesList(botID, pluginID, identifier string) (storagemodels.QuotesPluginQuotesList, error) {
	if q, ok := b.storage[botID][pluginID][identifier]; ok {
		if val, ok := q.(storagemodels.QuotesPluginQuotesList); ok {
			return val, nil
		}
		return storagemodels.QuotesPluginQuotesList{}, fmt.Errorf("Stored data is not a valid QuotesPluginQuotesList")
	}
	return storagemodels.QuotesPluginQuotesList{}, ErrNotFound
}
