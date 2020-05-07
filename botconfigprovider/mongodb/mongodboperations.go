package mongobotconfigprovider

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

func (b *MongoBotConfigProvider) createIndex(collection string, indexModel mongo.IndexModel) error {
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
func (b *MongoBotConfigProvider) checkBots() error {
	err := b.createIndex(COLLECTION_BOTS, mongo.IndexModel{
		Keys: bsonx.Doc{
			{Key: FIELD_ID, Value: bsonx.Int32(1)},
		},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return fmt.Errorf("Error creating MongoDB indices: %s", err)
	}

	return nil
}

// checkCollections checks if all collections needed exist and sets the correct indices
func (b *MongoBotConfigProvider) checkCollections() error {
	err := b.checkBots()
	if err != nil {
		return err
	}
	return nil
}
