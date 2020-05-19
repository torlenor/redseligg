package sqlitestorage

import (
	"fmt"
	"time"

	"github.com/torlenor/redseligg/storagemodels"
)

// GetQuotesPluginQuote returns a QuotesPluginQuote.
func (b *SQLiteStorage) GetQuotesPluginQuote(botID, pluginID, identifier string) (storagemodels.QuotesPluginQuote, error) {
	row, err := b.db.Query(
		fmt.Sprintf(`SELECT * FROM %s WHERE bot_id="%s" AND plugin_id="%s" and identifier="%s"`,
			tableQuotesPluginQuote,
			botID,
			pluginID,
			identifier))
	if err != nil {
		return storagemodels.QuotesPluginQuote{}, err
	}
	defer row.Close()

	quote := storagemodels.QuotesPluginQuote{}
	for row.Next() {
		var id int
		var gotBotID string
		var gotPluginID string
		var gotIdentifier string
		row.Scan(&id, &gotBotID, &gotPluginID, &gotIdentifier, &quote.Author, &quote.Added, &quote.AuthorID, &quote.ChannelID, &quote.Text)
	}

	return quote, nil
}

// GetQuotesPluginQuotesList returns a QuotesPluginQuotesList.
func (b *SQLiteStorage) GetQuotesPluginQuotesList(botID, pluginID, identifier string) (storagemodels.QuotesPluginQuotesList, error) {
	row, err := b.db.Query(
		fmt.Sprintf(`SELECT identifier FROM %s WHERE bot_id="%s" AND plugin_id="%s"`,
			tableQuotesPluginQuote,
			botID,
			pluginID))
	if err != nil {
		return storagemodels.QuotesPluginQuotesList{}, err
	}
	defer row.Close()

	quotes := storagemodels.QuotesPluginQuotesList{}
	for row.Next() {
		var gotIdentifier string
		row.Scan(&gotIdentifier)
		quotes.UUIDs = append(quotes.UUIDs, gotIdentifier)
	}

	return quotes, nil
}

// GetTimedMessagesPluginMessages returns a TimedMessagesPluginMessages.
func (b *SQLiteStorage) GetTimedMessagesPluginMessages(botID, pluginID, identifier string) (storagemodels.TimedMessagesPluginMessages, error) {
	row, err := b.db.Query(
		fmt.Sprintf(`SELECT * FROM %s WHERE bot_id="%s" AND plugin_id="%s" AND identifier="%s"`,
			tableTimedMessagesPluginMessage,
			botID,
			pluginID,
			identifier))
	if err != nil {
		return storagemodels.TimedMessagesPluginMessages{}, err
	}
	defer row.Close()

	messages := storagemodels.TimedMessagesPluginMessages{}
	for row.Next() {
		var id int
		var gotBotID string
		var gotPluginID string
		var gotIdentifier string
		var gotInterval int64
		message := storagemodels.TimedMessagesPluginMessage{}
		err := row.Scan(&id, &gotBotID, &gotPluginID, &gotIdentifier, &message.Text, &gotInterval, &message.ChannelID, &message.LastSent)
		if err != nil {
			b.log.Errorf("Error parsing SQLite select: %s", err)
			continue
		}
		message.Interval = time.Millisecond * time.Duration(gotInterval)
		messages.Messages = append(messages.Messages, message)
	}

	return messages, nil
}

// GetCustomCommandsPluginCommands returns CustomCommandsPluginCommands.
func (b *SQLiteStorage) GetCustomCommandsPluginCommands(botID, pluginID, identifier string) (storagemodels.CustomCommandsPluginCommands, error) {

	row, err := b.db.Query(
		fmt.Sprintf(`SELECT * FROM %s WHERE bot_id="%s" AND plugin_id="%s" AND identifier="%s"`,
			tableCustomCommandsPluginCommands,
			botID,
			pluginID,
			identifier))
	if err != nil {
		return storagemodels.CustomCommandsPluginCommands{}, err
	}
	defer row.Close()

	commands := storagemodels.CustomCommandsPluginCommands{}
	for row.Next() {
		var id int
		var gotBotID string
		var gotPluginID string
		var gotIdentifier string
		command := storagemodels.CustomCommandsPluginCommand{}
		err := row.Scan(&id, &gotBotID, &gotPluginID, &gotIdentifier, &command.Command, &command.Text, &command.ChannelID)
		if err != nil {
			b.log.Errorf("Error parsing SQLite select: %s", err)
			continue
		}
		commands.Commands = append(commands.Commands, command)
	}

	return commands, nil
}
