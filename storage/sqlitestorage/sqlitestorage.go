package sqlitestorage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3" // import sqlite

	"github.com/sirupsen/logrus"

	"github.com/torlenor/redseligg/botconfig"
	"github.com/torlenor/redseligg/logging"
)

var tableArchivePluginMessage = "archive_plugin_message"
var tableCustomCommandsPluginCommands = "custom_commands_plugin_commands"
var tableQuotesPluginQuote = "quotes_plugin_quote"
var tableTimedMessagesPluginMessage = "timed_messages_plugin_message"
var tableRssPluginSubscription = "rss_plugin_subscription"

// SQLiteStorage is a SQLite implementation of a storage.
type SQLiteStorage struct {
	log    *logrus.Entry
	dbFile string

	db *sql.DB
}

// New creates a new SQLiteStorage
func New(storageConfig botconfig.StorageConfig) (*SQLiteStorage, error) {
	cfg, err := parseConfig(storageConfig)
	if err != nil {
		return nil, fmt.Errorf("Error parsing config %v: %s", storageConfig, err)
	}
	b := &SQLiteStorage{
		log:    logging.Get("SQLite Storage Backend"),
		dbFile: cfg.DBFile,
	}

	return b, nil
}

// Connect to the SQLite DB or create it if it does not exist
func (s *SQLiteStorage) Connect() error {
	if _, err := os.Stat(s.dbFile); err != nil {
		s.log.Infof("Creating new SQLite Storage for DB=%s", s.dbFile)
		file, err := os.Create(s.dbFile)
		if err != nil {
			return err
		}
		file.Close()
		log.Println(s.dbFile, "created")
	}

	var err error
	s.db, err = sql.Open("sqlite3", s.dbFile)
	if err != nil {
		return err
	}

	err = s.createTables()
	if err != nil {
		return err
	}

	return nil
}
