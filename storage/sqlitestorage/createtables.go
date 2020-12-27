package sqlitestorage

import "fmt"

func (s *SQLiteStorage) createTableArchivePluginMessage() error {
	_, tableExists := s.db.Query("select * from " + tableArchivePluginMessage + ";")
	if tableExists == nil {
		s.log.Debugf("Table %s already exists. Skipping creation", tableArchivePluginMessage)
		return nil
	}

	creatTableSQL := fmt.Sprintf(`CREATE TABLE %s (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"bot_id" TEXT,
		"plugin_id" TEXT,
		"identifier" TEXT,
		"timestamp" DATETIME,
		"channel_id" TEXT,
		"channel" TEXT,
		"user_id" TEXT,
		"user_name" TEXT,
		"content" TEXT,
		"private" BOOL
	  );`, tableArchivePluginMessage)

	s.log.Printf("Creating %s table...", tableArchivePluginMessage)
	statement, err := s.db.Prepare(creatTableSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	s.log.Printf("%s table created", tableArchivePluginMessage)
	return nil
}

func (s *SQLiteStorage) createTableQuotesPluginQuote() error {
	_, tableExists := s.db.Query("select * from " + tableQuotesPluginQuote + ";")
	if tableExists == nil {
		s.log.Debugf("Table %s already exists. Skipping creation", tableQuotesPluginQuote)
		return nil
	}

	creatTableSQL := fmt.Sprintf(`CREATE TABLE %s (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"bot_id" TEXT,
		"plugin_id" TEXT,
		"identifier" TEXT,
		"author" TEXT,
		"added" DATETIME,
		"author_id" TEXT,
		"channel_id" TEXT,
		"text" TEXT
	  );`, tableQuotesPluginQuote)

	s.log.Printf("Creating %s table...", tableQuotesPluginQuote)
	statement, err := s.db.Prepare(creatTableSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	s.log.Printf("%s table created", tableQuotesPluginQuote)
	return nil
}

func (s *SQLiteStorage) createTableTimedMessagesPluginMessage() error {
	_, tableExists := s.db.Query("select * from " + tableTimedMessagesPluginMessage + ";")
	if tableExists == nil {
		s.log.Debugf("Table %s already exists. Skipping creation", tableTimedMessagesPluginMessage)
		return nil
	}

	creatTableSQL := fmt.Sprintf(`CREATE TABLE %s (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"bot_id" TEXT,
		"plugin_id" TEXT,
		"identifier" TEXT,
		"text" TEXT,
		"interval_ms" INTEGER,
		"channel_id" TEXT,
		"last_sent" DATETIME
	  );`, tableTimedMessagesPluginMessage)

	s.log.Printf("Creating %s table...", tableTimedMessagesPluginMessage)
	statement, err := s.db.Prepare(creatTableSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	s.log.Printf("%s table created", tableTimedMessagesPluginMessage)
	return nil
}

func (s *SQLiteStorage) createTableCustomCommandsPluginCommands() error {
	_, tableExists := s.db.Query("select * from " + tableCustomCommandsPluginCommands + ";")
	if tableExists == nil {
		s.log.Debugf("Table %s already exists. Skipping creation", tableCustomCommandsPluginCommands)
		return nil
	}

	creatTableSQL := fmt.Sprintf(`CREATE TABLE %s (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"bot_id" TEXT,
		"plugin_id" TEXT,
		"identifier" TEXT,
		"command" TEXT,
		"text" TEXT,
		"channel_id" TEXT
	  );`, tableCustomCommandsPluginCommands)

	s.log.Printf("Creating %s table...", tableCustomCommandsPluginCommands)
	statement, err := s.db.Prepare(creatTableSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	s.log.Printf("%s table created", tableCustomCommandsPluginCommands)
	return nil
}

func (s *SQLiteStorage) createTableRssPluginSubscription() error {
	_, tableExists := s.db.Query("select * from " + tableRssPluginSubscription + ";")
	if tableExists == nil {
		s.log.Debugf("Table %s already exists. Skipping creation", tableRssPluginSubscription)
		return nil
	}

	creatTableSQL := fmt.Sprintf(`CREATE TABLE %s (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
		"bot_id" TEXT,
		"plugin_id" TEXT,
		"identifier" TEXT,
		"link" TEXT,
		"channel_id" TEXT,
		"last_posted_pub_date" DATETIME
	  );`, tableRssPluginSubscription)

	s.log.Printf("Creating %s table...", tableRssPluginSubscription)
	statement, err := s.db.Prepare(creatTableSQL)
	if err != nil {
		return err
	}
	_, err = statement.Exec()
	if err != nil {
		return err
	}
	s.log.Printf("%s table created", tableRssPluginSubscription)
	return nil
}

func (s *SQLiteStorage) createTables() error {
	err := s.createTableArchivePluginMessage()
	if err != nil {
		return err
	}
	err = s.createTableQuotesPluginQuote()
	if err != nil {
		return err
	}
	err = s.createTableTimedMessagesPluginMessage()
	if err != nil {
		return err
	}
	err = s.createTableCustomCommandsPluginCommands()
	if err != nil {
		return err
	}
	err = s.createTableRssPluginSubscription()
	if err != nil {
		return err
	}
	return nil
}
