// Package storagemodels contains custom storage models used by plugins
package storagemodels

import (
	"fmt"
	"time"
)

// QuotesPluginQuote is a QuotesPlugin storage model.
// It is used to store a single quote.
type QuotesPluginQuote struct {
	Author string
	Added  time.Time

	AuthorID  string
	ChannelID string

	Text string
}

func (q QuotesPluginQuote) String() string {
	year, month, day := q.Added.Date()
	return fmt.Sprintf(`"%s" - %d-%d-%d, added by %s`, q.Text, year, month, day, q.Author)
}

// QuotesPluginQuotesList is a QuotesPlugin storage model.
// It is used to store a list of Quote UUIDs.
type QuotesPluginQuotesList struct {
	UUIDs []string
}
