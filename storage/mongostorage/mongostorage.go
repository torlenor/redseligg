package mongostorage

import (
	"context"
	"fmt"
	"reflect"

	"github.com/sirupsen/logrus"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/bsontype"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"git.abyle.org/redseligg/botorchestrator/botconfig"

	"github.com/torlenor/abylebotter/logging"
)

var collectionPluginStorage = "pluginstorage"

var fieldBotID = "bot_id"
var fieldPluginID = "plugin_id"
var fieldIdentifier = "identifier"

// MongoStorage is a MongoDB implementation of a storage.
type MongoStorage struct {
	log    *logrus.Entry
	dbName string

	client *mongo.Client
	db     *mongo.Database

	connected bool
}

// New creates a new MongoStorage
func New(storageConfig botconfig.StorageConfig) (*MongoStorage, error) {
	cfg, err := parseConfig(storageConfig)
	if err != nil {
		return nil, fmt.Errorf("Error parsing config %v: %s", storageConfig, err)
	}
	b := &MongoStorage{
		log:    logging.Get("MongoDB Storage Backend"),
		dbName: cfg.DB,
	}

	b.log.Infof("Creating new MongoDB Storage connection for URL=%s and DB=%s", cfg.URL, cfg.DB)

	tM := reflect.TypeOf(bson.M{})
	reg := bson.NewRegistryBuilder().RegisterTypeMapEntry(bsontype.EmbeddedDocument, tM).Build()

	clientOptions := options.Client().ApplyURI(cfg.URL).SetRegistry(reg)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		b.log.Errorf("Error creating new MongoDB client: %s", err)
		return nil, err
	}
	b.client = client

	return b, nil
}

// Connect initializes the MongoStorage by using the appropriate Mongo functions.
// In addition it checks its collections and creates indices if necessary.
// This method must be called before the Mongo Backend can be used.
func (b *MongoStorage) Connect() error {
	err := b.client.Connect(context.Background())
	if err != nil {
		b.log.Errorf("Error connecting to MongoDB: %s", err)
		return err
	}

	b.connected = true

	b.db = b.client.Database(b.dbName)

	err = b.checkCollections()
	if err != nil {
		b.log.Errorf("Error checking collections: %s", err)
		return err
	}

	return nil
}

// IsConnected indicates if there is a connection to a Mongo server
func (b *MongoStorage) IsConnected() bool {
	return b.connected
}
