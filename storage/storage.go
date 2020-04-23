package storage

// Backend is an interface which described a storage backend
type backend interface {
	StorePluginData(botID, pluginID, identifier string, data interface{}) error
	GetPluginData(botID, pluginID, identifier string) (interface{}, error)
}

// Storage handles data storage and acquisition to and from a backend
type Storage struct {
	backend backend

	botID string
}

// New creates a new Storage with the given backend and for the given bot ID
func New(backend backend, botID string) *Storage {
	return &Storage{
		backend: backend,

		botID: botID,
	}
}

// Store can be used to store arbitrary data.
func (s *Storage) Store(pluginID, identifier string, data interface{}) error {
	return s.backend.StorePluginData(s.botID, pluginID, identifier, data)
}

// Get can be used to retreive arbitrary data.
func (s *Storage) Get(pluginID, identifier string) (interface{}, error) {
	return s.backend.GetPluginData(s.botID, pluginID, identifier)
}
