// Package storagemodels contains custom storage models used by plugins
package storagemodels

import (
	"time"
)

// RssPluginSubscription is a RssPlugin storage model.
// It represents one stored RSS feed with its latest posted pubDate.
type RssPluginSubscription struct {
	Link string

	ChannelID  string
	Identifier string

	LastPostedPubDate time.Time
}

// RssPluginSubscriptions is a TimedMessagesPlugin storage model.
// It is used to store all RSS feeds with their last posted pubDate.
type RssPluginSubscriptions struct {
	Subscriptions []RssPluginSubscription
}
