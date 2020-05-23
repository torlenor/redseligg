// Package storagemodels contains custom storage models used by plugins
package storagemodels

import (
	"time"
)

// ArchivePluginMessage is a ArchivePlugin storage model.
// It is used to store a single message.
type ArchivePluginMessage struct {
	TImestamp time.Time

	ServerID string
	Server   string

	ChannelID string
	Channel   string

	UserID   string
	UserName string

	Content string

	IsPrivate bool
}
