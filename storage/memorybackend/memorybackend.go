package memorybackend

// MemoryBackend is a simple in-memory implementation of a storage backend.
// It serves mostly as an example on how to implement a real storage backend.
type MemoryBackend struct {
	storage botStorage
}

// New returns a new MemoryBackend
func New() *MemoryBackend {
	return &MemoryBackend{
		storage: make(botStorage),
	}
}
