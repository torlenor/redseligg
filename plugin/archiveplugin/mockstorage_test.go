package archiveplugin

import "github.com/torlenor/redseligg/storagemodels"

type storedMessage struct {
	BotID      string
	PluginID   string
	Identifier string

	Data storagemodels.ArchivePluginMessage
}

// MockStorage is a mock storage implementation and can be used for testing.
type MockStorage struct {
	StoredMessages []storedMessage

	ErrorToReturn error
}

// Reset the MockStorage
func (b *MockStorage) Reset() {
	b.StoredMessages = []storedMessage{}
}

func (b *MockStorage) StoreArchivePluginMessage(botID, pluginID, identifier string, data storagemodels.ArchivePluginMessage) error {
	b.StoredMessages = append(b.StoredMessages, storedMessage{
		BotID:      botID,
		PluginID:   pluginID,
		Identifier: identifier,
		Data:       data,
	})
	return b.ErrorToReturn
}
