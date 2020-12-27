package sqlitestorage

import (
	"fmt"
)

// DeleteQuotesPluginQuote deletes a QuotesPluginQuote.
func (b *SQLiteStorage) DeleteQuotesPluginQuote(botID, pluginID, identifier string) error {
	deleteSQL := fmt.Sprintf(`DELETE FROM %s WHERE bot_id=? AND plugin_id=? AND identifier=?`, tableQuotesPluginQuote)
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

// DeleteRssPluginSubscription deletes a RssPluginSubscription.
func (b *SQLiteStorage) DeleteRssPluginSubscription(botID, pluginID, identifier string) error {
	deleteSQL := fmt.Sprintf(`DELETE FROM %s WHERE bot_id=? AND plugin_id=? AND identifier=?`, tableRssPluginSubscription)
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
