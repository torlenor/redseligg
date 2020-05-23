package mongostorage

import (
	"context"

	"github.com/torlenor/redseligg/storagemodels"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StoreQuotesPluginQuote takes a QuotesPluginQuote and stores it.
func (b *MongoStorage) StoreQuotesPluginQuote(botID, pluginID, identifier string, data storagemodels.QuotesPluginQuote) error {
	c := b.db.Collection(collectionPluginStorage)
	filter := bson.M{fieldBotID: botID, fieldPluginID: pluginID, fieldIdentifier: identifier}
	_, err := c.ReplaceOne(context.Background(), filter, quotesPluginQuoteData{BotID: botID, PluginID: pluginID, Identifier: identifier, Data: data}, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

// StoreQuotesPluginQuotesList takes a QuotesPluginQuotesList and stores it.
func (b *MongoStorage) StoreQuotesPluginQuotesList(botID, pluginID, identifier string, data storagemodels.QuotesPluginQuotesList) error {
	c := b.db.Collection(collectionPluginStorage)
	filter := bson.M{fieldBotID: botID, fieldPluginID: pluginID, fieldIdentifier: identifier}
	_, err := c.ReplaceOne(context.Background(), filter, quotesPluginQuotesListData{BotID: botID, PluginID: pluginID, Identifier: identifier, Data: data}, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

// StoreTimedMessagesPluginMessages stores data for TimedMessagesPlugin.
func (b *MongoStorage) StoreTimedMessagesPluginMessages(botID, pluginID, identifier string, data storagemodels.TimedMessagesPluginMessages) error {
	c := b.db.Collection(collectionPluginStorage)
	filter := bson.M{fieldBotID: botID, fieldPluginID: pluginID, fieldIdentifier: identifier}
	_, err := c.ReplaceOne(context.Background(), filter, timedMessagesPluginMessagesData{BotID: botID, PluginID: pluginID, Identifier: identifier, Data: data}, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

// StoreCustomCommandsPluginCommands stores data for CustomCommandsPlugin.
func (b *MongoStorage) StoreCustomCommandsPluginCommands(botID, pluginID, identifier string, data storagemodels.CustomCommandsPluginCommands) error {
	c := b.db.Collection(collectionPluginStorage)
	filter := bson.M{fieldBotID: botID, fieldPluginID: pluginID, fieldIdentifier: identifier}
	_, err := c.ReplaceOne(context.Background(), filter, customCommandsPluginCommandsData{BotID: botID, PluginID: pluginID, Identifier: identifier, Data: data}, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

// StoreArchivePluginMessage stores data for ArchivePlugin.
func (b *MongoStorage) StoreArchivePluginMessage(botID, pluginID, identifier string, data storagemodels.ArchivePluginMessage) error {
	c := b.db.Collection(collectionPluginStorage)
	filter := bson.M{fieldBotID: botID, fieldPluginID: pluginID, fieldIdentifier: identifier}
	_, err := c.ReplaceOne(context.Background(), filter,
		archivePluginMessageData{BotID: botID, PluginID: pluginID, Identifier: identifier, Data: data},
		options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}
