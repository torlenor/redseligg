package rssplugin

import (
	"fmt"

	"github.com/torlenor/redseligg/storagemodels"
)

type storedData struct {
	BotID      string
	PluginID   string
	Identifier string
	Data       storagemodels.RssPluginSubscription
}

type retrievedData struct {
	BotID      string
	PluginID   string
	Identifier string
}

// MockStorage is a mock storage implementation and can be used for testing
type MockStorage struct {
	StoredSubscriptions []storedData

	LastRetrieved retrievedData

	DataToReturn  storagemodels.RssPluginSubscriptions
	ErrorToReturn error
}

// Reset the MockStorage
func (b *MockStorage) Reset() {
	b.StoredSubscriptions = []storedData{}
	b.LastRetrieved = retrievedData{}
}

func (b *MockStorage) StoreRssPluginSubscription(botID, pluginID, identifier string, data storagemodels.RssPluginSubscription) error {
	b.StoredSubscriptions = append(b.StoredSubscriptions, storedData{
		PluginID:   pluginID,
		Identifier: identifier,
		Data:       data,
	})

	fmt.Printf("STORED %s\n", data.Link)
	fmt.Printf("Length = %d\n", len(b.StoredSubscriptions))

	return b.ErrorToReturn
}

func (b *MockStorage) GetRssPluginSubscriptions(botID, pluginID string) (storagemodels.RssPluginSubscriptions, error) {
	b.LastRetrieved.BotID = botID
	b.LastRetrieved.PluginID = pluginID

	return b.DataToReturn, b.ErrorToReturn
}

// DeleteRssPluginSubscription takes a RssPluginSubscription and updates it.
func (b *MockStorage) DeleteRssPluginSubscription(botID, pluginID, identifier string) error {
	// TODO: Test DeleteRssPluginSubscription behavior of RssPlugin

	return nil
}

// UpdateRssPluginSubscription takes a RssPluginSubscription and updates it.
func (b *MockStorage) UpdateRssPluginSubscription(botID, pluginID, identifier string, data storagemodels.RssPluginSubscription) error {
	// TODO: Test UpdateRssPluginSubscription behavior of RssPlugin

	return nil
}
