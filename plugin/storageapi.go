package plugin

// The StorageAPI can be used to store and get data from storage.
//
// A valid Bot has to provide a storage implmenetation with all the defined functions below.
//
// Plugins obtain access to this by embedding AbyleBotterPlugin.
type StorageAPI interface {
	// Store can be used to store arbitrary data.
	Store(pluginID, identifier string, data interface{}) error

	// Get can be used to retreive arbitrary data.
	Get(pluginID, identifier string) (interface{}, error)
}
