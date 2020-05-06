// Package storagemodels contains custom storage models used by plugins
package storagemodels

import (
	"time"
)

// TimedMessagesPluginMessage is a TimedMessagesPlugin storage model.
// It represents one timed messages with its interval.
type TimedMessagesPluginMessage struct {
	Text     string
	Interval time.Duration

	ChannelID string
}

// TimedMessagesPluginMessages is a TimedMessagesPlugin storage model.
// It is used to store all timed messages with their intervals.
type TimedMessagesPluginMessages struct {
	Messages []TimedMessagesPluginMessage
}
