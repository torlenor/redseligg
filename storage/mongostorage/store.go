package mongostorage

import (
	"context"

	"github.com/torlenor/abylebotter/storagemodels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (b *MongoStorage) StoreQuotesPluginQuote(botID, pluginID, identifier string, data storagemodels.QuotesPluginQuote) error {
	c := b.db.Collection(collectionPluginStorage)
	filter := bson.M{fieldBotID: botID, fieldPluginID: pluginID, fieldIdentifier: identifier}
	_, err := c.ReplaceOne(context.Background(), filter, quotesPluginQuoteData{BotID: botID, PluginID: pluginID, Identifier: identifier, Data: data}, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

func (b *MongoStorage) StoreQuotesPluginQuotesList(botID, pluginID, identifier string, data storagemodels.QuotesPluginQuotesList) error {
	c := b.db.Collection(collectionPluginStorage)
	filter := bson.M{fieldBotID: botID, fieldPluginID: pluginID, fieldIdentifier: identifier}
	_, err := c.ReplaceOne(context.Background(), filter, quotesPluginQuotesListData{BotID: botID, PluginID: pluginID, Identifier: identifier, Data: data}, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}
