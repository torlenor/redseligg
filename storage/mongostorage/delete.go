package mongostorage

import (
	"context"
	"fmt"

	"github.com/torlenor/abylebotter/storage"
	"go.mongodb.org/mongo-driver/bson"
)

// DeleteQuotesPluginQuote deletes a QuotesPluginQuote.
func (b *MongoStorage) DeleteQuotesPluginQuote(botID, pluginID, identifier string) error {
	if !b.IsConnected() {
		return fmt.Errorf("Not connected to MongoDB")
	}

	c := b.db.Collection(collectionPluginStorage)

	filter := bson.M{fieldBotID: botID, fieldPluginID: pluginID, fieldIdentifier: identifier}
	res, err := c.DeleteOne(context.Background(), filter)
	if err != nil {
		return fmt.Errorf("Error occurred deleting quote: %s", err)
	}
	if res.DeletedCount == 0 {
		return storage.ErrNotFound
	}

	return nil
}
