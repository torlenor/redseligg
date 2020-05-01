package memorystorage

// MemoryStorage is a simple in-memory implementation of a storage backend.
// It serves mostly as an example on how to implement a real storage backend.
type MemoryStorage struct {
	storage botStorage
}

// New returns a new MemoryStorage
func New() *MemoryStorage {
	return &MemoryStorage{
		storage: make(botStorage),
	}
}
