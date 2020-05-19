package sqlitestorage

import (
	"fmt"

	"github.com/torlenor/redseligg/storagemodels"
)

// StoreQuotesPluginQuote takes a QuotesPluginQuote and stores it.
func (b *SQLiteStorage) StoreQuotesPluginQuote(botID, pluginID, identifier string, data storagemodels.QuotesPluginQuote) error {
	insertSQL := fmt.Sprintf(`INSERT INTO %s(bot_id, plugin_id, identifier, author, added, author_id, channel_id, text) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`, tableQuotesPluginQuote)
	statement, err := b.db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("Could not prepare sql statement: %s", err)
	}
	_, err = statement.Exec(botID, pluginID, identifier, data.Author, data.Added, data.AuthorID, data.ChannelID, data.Text)
	if err != nil {
		return fmt.Errorf("Could not insert data: %s", err)
	}

	return nil
}

// StoreQuotesPluginQuotesList takes a QuotesPluginQuotesList and stores it.
func (b *SQLiteStorage) StoreQuotesPluginQuotesList(botID, pluginID, identifier string, data storagemodels.QuotesPluginQuotesList) error {
	// Not needed
	return nil
}

func (b *SQLiteStorage) clearEntries(table, botID, pluginID, identifier string) error {
	deleteSQL := fmt.Sprintf(`DELETE FROM %s WHERE bot_id=? AND plugin_id=? AND identifier=?`, table)
	statement, err := b.db.Prepare(deleteSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec(botID, pluginID, identifier)
	if err != nil {
		return err
	}
	return nil
}

// StoreTimedMessagesPluginMessages stores data for TimedMessagesPlugin
func (b *SQLiteStorage) StoreTimedMessagesPluginMessages(botID, pluginID, identifier string, data storagemodels.TimedMessagesPluginMessages) error {
	err := b.clearEntries(tableTimedMessagesPluginMessage, botID, pluginID, identifier)
	if err != nil {
		return fmt.Errorf("Could not clean up table: %s", err)
	}

	insertSQL := fmt.Sprintf(`INSERT INTO %s(bot_id, plugin_id, identifier, text, interval_ms, channel_id, last_sent) VALUES (?, ?, ?, ?, ?, ?, ?)`, tableTimedMessagesPluginMessage)
	statement, err := b.db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("Could not prepare sql statement: %s", err)
	}

	for _, message := range data.Messages {
		_, err = statement.Exec(botID, pluginID, identifier, message.Text, message.Interval.Milliseconds(), message.ChannelID, message.LastSent)
		if err != nil {
			return fmt.Errorf("Could not insert data: %s", err)
		}
	}

	return nil
}

// StoreCustomCommandsPluginCommands stores data for CustomCommandsPlugin
func (b *SQLiteStorage) StoreCustomCommandsPluginCommands(botID, pluginID, identifier string, data storagemodels.CustomCommandsPluginCommands) error {
	err := b.clearEntries(tableCustomCommandsPluginCommands, botID, pluginID, identifier)
	if err != nil {
		return fmt.Errorf("Could not clean up table: %s", err)
	}

	insertSQL := fmt.Sprintf(`INSERT INTO %s(bot_id, plugin_id, identifier, command, text, channel_id) VALUES (?, ?, ?, ?, ?, ?)`, tableCustomCommandsPluginCommands)
	statement, err := b.db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("Could not prepare sql statement: %s", err)
	}

	for _, command := range data.Commands {
		_, err = statement.Exec(botID, pluginID, identifier, command.Command, command.Text, command.ChannelID)
		if err != nil {
			return fmt.Errorf("Could not insert data: %s", err)
		}
	}

	return nil
}
