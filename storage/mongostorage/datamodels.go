package mongostorage

import "github.com/torlenor/abylebotter/storagemodels"

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
