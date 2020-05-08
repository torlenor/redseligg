package mongostorage

import "github.com/torlenor/redseligg/storagemodels"

type quotesPluginQuoteData struct {
	BotID      string `bson:"bot_id"`
	PluginID   string `bson:"plugin_id"`
	Identifier string `bson:"identifier"`

	Data storagemodels.QuotesPluginQuote `bson:"data"`
}

type quotesPluginQuotesListData struct {
	BotID      string `bson:"bot_id"`
	PluginID   string `bson:"plugin_id"`
	Identifier string `bson:"identifier"`

	Data storagemodels.QuotesPluginQuotesList `bson:"data"`
}

type timedMessagesPluginMessagesData struct {
	BotID      string `bson:"bot_id"`
	PluginID   string `bson:"plugin_id"`
	Identifier string `bson:"identifier"`

	Data storagemodels.TimedMessagesPluginMessages `bson:"data"`
}

type customCommandsPluginCommandsData struct {
	BotID      string `bson:"bot_id"`
	PluginID   string `bson:"plugin_id"`
	Identifier string `bson:"identifier"`

	Data storagemodels.CustomCommandsPluginCommands `bson:"data"`
}
