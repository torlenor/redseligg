package mongostorage

import (
	"context"
	"fmt"

	"github.com/torlenor/redseligg/storage"
	"github.com/torlenor/redseligg/storagemodels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
		return storagemodels.QuotesPluginQuote{}, storage.ErrNotFound
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
		return storagemodels.QuotesPluginQuotesList{}, storage.ErrNotFound
	} else if err != nil {
		return storagemodels.QuotesPluginQuotesList{}, fmt.Errorf("Error in finding the bot config with id %s: %s", botID, err)
	}

	return data.Data, nil
}

// GetTimedMessagesPluginMessages returns a TimedMessagesPluginMessages.
func (b *MongoStorage) GetTimedMessagesPluginMessages(botID, pluginID, identifier string) (storagemodels.TimedMessagesPluginMessages, error) {
	if !b.IsConnected() {
		return storagemodels.TimedMessagesPluginMessages{}, fmt.Errorf("Not connected to MongoDB")
	}

	c := b.db.Collection(collectionPluginStorage)

	filter := bson.M{fieldBotID: botID, fieldPluginID: pluginID, fieldIdentifier: identifier}
	var data timedMessagesPluginMessagesData
	err := c.FindOne(context.Background(), filter).Decode(&data)
	if err == mongo.ErrNoDocuments {
		return storagemodels.TimedMessagesPluginMessages{}, storage.ErrNotFound
	} else if err != nil {
		return storagemodels.TimedMessagesPluginMessages{}, fmt.Errorf("Error in finding the bot config with id %s: %s", botID, err)
	}

	return data.Data, nil
}

// GetCustomCommandsPluginCommands returns CustomCommandsPluginCommands.
func (b *MongoStorage) GetCustomCommandsPluginCommands(botID, pluginID, identifier string) (storagemodels.CustomCommandsPluginCommands, error) {
	if !b.IsConnected() {
		return storagemodels.CustomCommandsPluginCommands{}, fmt.Errorf("Not connected to MongoDB")
	}

	c := b.db.Collection(collectionPluginStorage)

	filter := bson.M{fieldBotID: botID, fieldPluginID: pluginID, fieldIdentifier: identifier}
	var data customCommandsPluginCommandsData
	err := c.FindOne(context.Background(), filter).Decode(&data)
	if err == mongo.ErrNoDocuments {
		return storagemodels.CustomCommandsPluginCommands{}, storage.ErrNotFound
	} else if err != nil {
		return storagemodels.CustomCommandsPluginCommands{}, fmt.Errorf("Error in finding the bot config with id %s: %s", botID, err)
	}

	return data.Data, nil
}
