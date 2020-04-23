package plugin

// The MockStorageAPI can be used for testing Plugins by providing helper functions.
// It mimics all the functions a real bot would, but in addition provides helper functions
// for unit tests.

type storedData struct {
	PluginID   string
	Identifier string
	Data       interface{}
}

type retrievedData struct {
	PluginID   string
	Identifier string
}

// MockStorageAPI is a mock storage implementation and can be used for testing
type MockStorageAPI struct {
	Stored []storedData

	LastRetrieved retrievedData

	DataToReturn  storedData
	ErrorToReturn error
}

// Reset the MockStorageAPI
func (b *MockStorageAPI) Reset() {
}

// Store can be used to store arbitrary data.
func (b *MockStorageAPI) Store(pluginID, identifier string, data interface{}) error {
	b.Stored = append(b.Stored, storedData{
		PluginID:   pluginID,
		Identifier: identifier,
		Data:       data,
	})
	return b.ErrorToReturn
}

// Get can be used to retreive arbitrary data.
func (b *MockStorageAPI) Get(pluginID, identifier string) (interface{}, error) {
	b.LastRetrieved.PluginID = pluginID
	b.LastRetrieved.Identifier = identifier

	return b.DataToReturn.Data, b.ErrorToReturn
}
