package sqlitestorage

import (
	"fmt"

	"github.com/torlenor/redseligg/botconfig"
)

type config struct {
	DBFile string
}

func parseConfig(c botconfig.StorageConfig) (config, error) {
	if c.Type != "sqlite" {
		return config{}, fmt.Errorf("Not a SQLite storage config")
	}

	var dbFile string

	var ok bool
	if dbFile, ok = c.Config["database"].(string); !ok {
		return config{}, fmt.Errorf("Database file not defined in config")
	}

	cfg := config{
		DBFile: dbFile,
	}

	return cfg, nil
}
