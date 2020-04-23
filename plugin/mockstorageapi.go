package plugin

import (
	"fmt"
)

// The MockStorageAPI can be used for testing Plugins by providing helper functions.
// It mimics all the functions a real bot would, but in addition provides helper functions
// for unit tests.
type MockStorageAPI struct{}

// Reset the MockStorageAPI
func (b *MockStorageAPI) Reset() {
}

// Store can be used to store arbitrary data.
func (b *MockStorageAPI) Store(pluginID, identifier string, data interface{}) error {
	return fmt.Errorf("Not implemented")
}

// Get can be used to retreive arbitrary data.
func (b *MockStorageAPI) Get(pluginID, identifier string) (interface{}, error) {
	return nil, fmt.Errorf("Not implemented")
}
