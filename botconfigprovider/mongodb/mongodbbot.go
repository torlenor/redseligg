package mongobotconfigprovider

import (
	"context"
	"fmt"

	"github.com/torlenor/redseligg/botconfig"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var FIELD_ID = "id"
var FIELD_ENABLED = "enabled"
var COLLECTION_BOTS = "bots"

// GetBotConfig retrieves the bot configuration data for the given id
func (b *MongoBotConfigProvider) GetBotConfig(botID string) (botconfig.BotConfig, error) {
	if !b.IsConnected() {
		return botconfig.BotConfig{}, fmt.Errorf("Not connected to MongoDB")
	}

	c := b.db.Collection(COLLECTION_BOTS)

	filter := bson.D{{FIELD_ID, botID}}
	botConfig := botconfig.BotConfig{}
	err := c.FindOne(context.Background(), filter).Decode(&botConfig)
	if err != nil {
		return botconfig.BotConfig{}, fmt.Errorf("Error in finding the bot config with id %s: %s", botID, err)
	}

	return botConfig, nil
}

// GetBotConfigCount returns the number of stored bot configurations
func (b *MongoBotConfigProvider) GetBotConfigCount() (uint64, error) {
	if !b.IsConnected() {
		return 0, fmt.Errorf("Not connected to MongoDB")
	}

	c := b.db.Collection(COLLECTION_BOTS)

	count, err := c.CountDocuments(
		context.Background(),
		bson.D{{}},
	)
	if err != nil {
		return uint64(0), fmt.Errorf("Error in counting bot configs: %s", err)
	}

	return uint64(count), nil
}

// StoreBotConfig stores new match data
func (b *MongoBotConfigProvider) StoreBotConfig(data *botconfig.BotConfig) error {
	if !b.IsConnected() {
		return fmt.Errorf("Not connected to MongoDB")
	}

	c := b.db.Collection(COLLECTION_BOTS)
	filter := bson.M{FIELD_ID: data.BotID}
	_, err := c.ReplaceOne(context.Background(), filter, data, options.Replace().SetUpsert(true))
	if err != nil {
		return err
	}

	return nil
}

// GetAllEnabledBotIDs returns only bot IDs for bots which are enabled
// It returns nil if something went wrong
func (b *MongoBotConfigProvider) GetAllEnabledBotIDs() (botIDs []string) {
	c := b.db.Collection(COLLECTION_BOTS)

	filter := bson.D{{FIELD_ENABLED, true}}

	cur, err := c.Find(
		context.Background(),
		filter,
	)
	if err != nil {
		b.log.Errorf("Find error: %s", err)
		return nil
	}

	defer cur.Close(context.Background())

	for cur.Next(nil) {
		botConfig := botconfig.BotConfig{}
		err := cur.Decode(&botConfig)
		if err != nil {
			b.log.Warnln("Decode error ", err)
			continue
		}
		botIDs = append(botIDs, botConfig.BotID)
	}

	if err := cur.Err(); err != nil {
		b.log.Warnln("Cursor error ", err)
	}

	return botIDs
}
