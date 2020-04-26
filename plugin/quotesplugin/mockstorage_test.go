package quotesplugin

import "github.com/torlenor/abylebotter/storagemodels"

type storedQuoteData struct {
	BotID      string
	PluginID   string
	Identifier string
	Data       storagemodels.QuotesPluginQuote
}

type storedQuotesListData struct {
	BotID      string
	PluginID   string
	Identifier string
	Data       storagemodels.QuotesPluginQuotesList
}

type retrievedData struct {
	BotID      string
	PluginID   string
	Identifier string
}

// MockStorage is a mock storage implementation and can be used for testing
type MockStorage struct {
	StoredQuotes     []storedQuoteData
	StoredQuotesList []storedQuotesListData

	LastRetrieved retrievedData

	QuoteDataToReturn      storagemodels.QuotesPluginQuote
	QuotesListDataToReturn storagemodels.QuotesPluginQuotesList
	ErrorToReturn          error
}

// Reset the MockStorage
func (b *MockStorage) Reset() {}

func (b *MockStorage) StoreQuotesPluginQuote(botID, pluginID, identifier string, data storagemodels.QuotesPluginQuote) error {
	b.StoredQuotes = append(b.StoredQuotes, storedQuoteData{
		PluginID:   pluginID,
		Identifier: identifier,
		Data:       data,
	})
	return b.ErrorToReturn
}

func (b *MockStorage) StoreQuotesPluginQuotesList(botID, pluginID, identifier string, data storagemodels.QuotesPluginQuotesList) error {
	b.StoredQuotesList = append(b.StoredQuotesList, storedQuotesListData{
		PluginID:   pluginID,
		Identifier: identifier,
		Data:       data,
	})
	return b.ErrorToReturn
}

func (b *MockStorage) GetQuotesPluginQuote(botID, pluginID, identifier string) (storagemodels.QuotesPluginQuote, error) {
	b.LastRetrieved.BotID = botID
	b.LastRetrieved.PluginID = pluginID
	b.LastRetrieved.Identifier = identifier

	return b.QuoteDataToReturn, b.ErrorToReturn
}

func (b *MockStorage) GetQuotesPluginQuotesList(botID, pluginID, identifier string) (storagemodels.QuotesPluginQuotesList, error) {
	b.LastRetrieved.BotID = botID
	b.LastRetrieved.PluginID = pluginID
	b.LastRetrieved.Identifier = identifier

	return b.QuotesListDataToReturn, b.ErrorToReturn
}
