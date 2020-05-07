package mongobotconfigprovider

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/torlenor/abylebotter/logging"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoBotConfigProvider represents a config provider using a MongoDB
type MongoBotConfigProvider struct {
	log    *logrus.Entry
	dbName string

	client *mongo.Client
	db     *mongo.Database

	connected bool
}

// NewBackend creates a new Mongo Backend
func NewBackend(url string, db string) (*MongoBotConfigProvider, error) {
	b := &MongoBotConfigProvider{
		log:    logging.Get("MongoDB Storage Backend"),
		dbName: db,
	}

	clientOptions := options.Client().ApplyURI(url)
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		b.log.Errorf("Error creating new MongoDB client: %s", err)
		return nil, err
	}
	b.client = client

	return b, nil
}

// Connect initializes the MongoDBBackend Client by using the appropriate Mongo functions.
// In addition it checks its collections and creates indices if necessary.
// This method must be called before the Mongo Backend can be used.
func (b *MongoBotConfigProvider) Connect() error {
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

// IsConnected indicates if the MongoBotConfigProvider actually has a connection to a Mongo server
func (b *MongoBotConfigProvider) IsConnected() bool {
	return b.connected
}
