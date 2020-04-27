package mongostorage

import (
	"context"
	"errors"
	"fmt"

	"github.com/torlenor/abylebotter/storagemodels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// ErrNotFound is returned when no data could be found
var ErrNotFound = errors.New("MongoStorage: Could not find requested data")

// GetQuotesPluginQuote returns a QuotesPluginQuote.
func (b *MongoStorage) GetQuotesPluginQuote(botID, pluginID, identifier string) (storagemodels.QuotesPluginQuote, error) {
	if !b.IsConnected() {
		return storagemodels.QuotesPluginQuote{}, fmt.Errorf("Not connected to MongoDB")
	}

	c := b.db.Collection(collectionPluginStorage)

	filter := bson.M{fieldBotID: botID, fieldPluginID: pluginID, fieldIdentifier: identifier}
	var data quotesPluginQuoteData
	err := c.FindOne(context.Background(), filter).Decode(&data)
	if err == mongo.ErrNoDocuments {
		return storagemodels.QuotesPluginQuote{}, fmt.Errorf("Could not find data for botID=%s, pluginID=%s, identifier=%s: %s", botID, pluginID, identifier, err)
	} else if err != nil {
		return storagemodels.QuotesPluginQuote{}, fmt.Errorf("Error in finding the bot config with id %s: %s", botID, err)
	}

	return data.Data, nil
}

// GetQuotesPluginQuotesList returns a QuotesPluginQuotesList.
func (b *MongoStorage) GetQuotesPluginQuotesList(botID, pluginID, identifier string) (storagemodels.QuotesPluginQuotesList, error) {
	if !b.IsConnected() {
		return storagemodels.QuotesPluginQuotesList{}, fmt.Errorf("Not connected to MongoDB")
	}

	c := b.db.Collection(collectionPluginStorage)

	filter := bson.M{fieldBotID: botID, fieldPluginID: pluginID, fieldIdentifier: identifier}
	var data quotesPluginQuotesListData
	err := c.FindOne(context.Background(), filter).Decode(&data)
	if err == mongo.ErrNoDocuments {
		return storagemodels.QuotesPluginQuotesList{}, fmt.Errorf("Could not find data for botID=%s, pluginID=%s, identifier=%s: %s", botID, pluginID, identifier, err)
	} else if err != nil {
		return storagemodels.QuotesPluginQuotesList{}, fmt.Errorf("Error in finding the bot config with id %s: %s", botID, err)
	}

	return data.Data, nil
}
