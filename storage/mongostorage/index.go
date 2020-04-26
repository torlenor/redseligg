package mongostorage

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func (b *MongoStorage) createIndex(collection string, indexModel mongo.IndexModel) error {
	indexView := b.db.Collection(collection).Indexes()
	_, err := indexView.CreateOne(
		context.Background(),
		indexModel,
	)
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	return nil
}

// checkMatches checks the matches collection and sets the correct indices
func (b *MongoStorage) checkPluginStorage() error {
	err := b.createIndex(collectionPluginStorage, mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: fieldBotID, Value: bsonx.Int32(1)},
			{Key: fieldPluginID, Value: bsonx.Int32(1)},
			{Key: fieldIdentifier, Value: bsonx.Int32(1)},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	return nil
}

// checkCollections checks if all collections needed exist and sets the correct indices
func (b *MongoStorage) checkCollections() error {
	err := b.checkPluginStorage()
	if err != nil {
		return err
	}
	return nil
}
