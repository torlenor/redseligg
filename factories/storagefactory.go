package factories

import (
	"fmt"

	"github.com/torlenor/redseligg/botconfig"

	"github.com/torlenor/redseligg/logging"
	"github.com/torlenor/redseligg/storage"
	"github.com/torlenor/redseligg/storage/memorystorage"
	"github.com/torlenor/redseligg/storage/mongostorage"
)

var (
	logStorageFactory = logging.Get("StorageFactory")
)

// StorageFactory can be used to generate storage backends
type StorageFactory struct{}

// CreateBackend creates a new storage backend with the provided configuration
func (f *StorageFactory) CreateBackend(storageConfig botconfig.StorageConfig) (storage.Storage, error) {
	var b storage.Storage

	switch storageConfig.Type {
	case "memory":
		logStorageFactory.Tracef("Creating Memory storage")
		b = memorystorage.New()
	case "mongodb":
		fallthrough
	case "mongo":
		logStorageFactory.Tracef("Creating Mongo storage")
		m, err := mongostorage.New(storageConfig)
		if err != nil {
			return nil, err
		}
		err = m.Connect()
		if err != nil {
			return nil, err
		}
		b = m
	case "":
		return nil, fmt.Errorf("No storage defined in config")
	default:
		return nil, fmt.Errorf("Unknown storage type %s", storageConfig.Type)
	}

	return b, nil
}
