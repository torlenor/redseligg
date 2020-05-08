package customcommandsplugin

import "github.com/torlenor/redseligg/storagemodels"

type storedData struct {
	BotID      string
	PluginID   string
	Identifier string
	Data       storagemodels.CustomCommandsPluginCommands
}

type retrievedData struct {
	BotID      string
	PluginID   string
	Identifier string
}

// MockStorage is a mock storage implementation and can be used for testing
type MockStorage struct {
	StoredData storedData

	LastRetrieved retrievedData

	DataToReturn  storagemodels.CustomCommandsPluginCommands
	ErrorToReturn error
}

// Reset the MockStorage
func (b *MockStorage) Reset() {
	b.StoredData = storedData{}
	b.LastRetrieved = retrievedData{}
}

func (b *MockStorage) StoreCustomCommandsPluginCommands(botID, pluginID, identifier string, data storagemodels.CustomCommandsPluginCommands) error {
	b.StoredData = storedData{
		PluginID:   pluginID,
		Identifier: identifier,
		Data:       data,
	}
	return b.ErrorToReturn
}

func (b *MockStorage) GetCustomCommandsPluginCommands(botID, pluginID, identifier string) (storagemodels.CustomCommandsPluginCommands, error) {
	b.LastRetrieved.BotID = botID
	b.LastRetrieved.PluginID = pluginID
	b.LastRetrieved.Identifier = identifier

	return b.DataToReturn, b.ErrorToReturn
}
