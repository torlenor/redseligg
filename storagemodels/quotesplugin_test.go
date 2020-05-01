// Package storagemodels contains custom storage models used by plugins
package storagemodels

import (
	"fmt"
	"testing"
	"time"
)

func TestQuotesPluginQuote_String(t *testing.T) {
	now := time.Now()
	year, month, day := now.Date()

	quote := QuotesPluginQuote{
		Author: "AUTHOR",
		Added:  now,

		AuthorID:  "AID",
		ChannelID: "CID",

		Text: "SOME Quote",
	}

	tests := []struct {
		name string
		q    QuotesPluginQuote
		want string
	}{
		{
			name: "Regular quote",
			q:    quote,
			want: fmt.Sprintf(`"%s" - %d-%d-%d, added by %s`, quote.Text, year, month, day, quote.Author),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.q.String(); got != tt.want {
				t.Errorf("QuotesPluginQuote.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
