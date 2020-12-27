package sqlitestorage

import (
	"fmt"

	"github.com/torlenor/redseligg/storagemodels"
)

// UpdateRssPluginSubscription takes a RssPluginSubscription and updates it.
func (b *SQLiteStorage) UpdateRssPluginSubscription(botID, pluginID, identifier string, data storagemodels.RssPluginSubscription) error {
	insertSQL := fmt.Sprintf(`UPDATE %s SET link=?, channel_id=?, last_posted_pub_date=? WHERE bot_id=? AND plugin_id=? AND identifier=?`, tableRssPluginSubscription)
	statement, err := b.db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("Could not prepare sql statement: %s", err)
	}
	_, err = statement.Exec(data.Link, data.ChannelID, data.LastPostedPubDate, botID, pluginID, identifier)
	if err != nil {
		return fmt.Errorf("Could not insert data: %s", err)
	}

	return nil
}
