package storagemodels

type QuotesPluginQuote struct {
	Author   string
	AuthorID string

	ChannelID string

	Text string
}

type QuotesPluginQuotesList struct {
	UUIDs []string
}
