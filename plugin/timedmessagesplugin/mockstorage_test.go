package timedmessagesplugin

import "github.com/torlenor/redseligg/storagemodels"

type storedData struct {
	BotID      string
	PluginID   string
	Identifier string
	Data       storagemodels.TimedMessagesPluginMessages
}

type retrievedData struct {
	BotID      string
	PluginID   string
	Identifier string
}

// MockStorage is a mock storage implementation and can be used for testing
type MockStorage struct {
	StoredMessages storedData

	LastRetrieved retrievedData

	DataToReturn  storagemodels.TimedMessagesPluginMessages
	ErrorToReturn error
}

// Reset the MockStorage
func (b *MockStorage) Reset() {
	b.StoredMessages = storedData{}
	b.LastRetrieved = retrievedData{}
}

func (b *MockStorage) StoreTimedMessagesPluginMessages(botID, pluginID, identifier string, data storagemodels.TimedMessagesPluginMessages) error {
	b.StoredMessages = storedData{
		PluginID:   pluginID,
		Identifier: identifier,
		Data:       data,
	}
	return b.ErrorToReturn
}

func (b *MockStorage) GetTimedMessagesPluginMessages(botID, pluginID, identifier string) (storagemodels.TimedMessagesPluginMessages, error) {
	b.LastRetrieved.BotID = botID
	b.LastRetrieved.PluginID = pluginID
	b.LastRetrieved.Identifier = identifier

	return b.DataToReturn, b.ErrorToReturn
}
