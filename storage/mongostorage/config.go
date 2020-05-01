package mongostorage

import (
	"fmt"

	"git.abyle.org/redseligg/botorchestrator/botconfig"
)

type config struct {
	URL string
	DB  string
}

func parseConfig(c botconfig.StorageConfig) (config, error) {
	if c.Type != "mongo" {
		return config{}, fmt.Errorf("Not a Mongo config")
	}

	var url string
	var db string

	var ok bool
	if url, ok = c.Config["url"].(string); !ok {
		return config{}, fmt.Errorf("URL not defined in config")
	}
	if db, ok = c.Config["database"].(string); !ok {
		return config{}, fmt.Errorf("Database not defined in config")
	}

	cfg := config{
		URL: url,
		DB:  db,
	}

	return cfg, nil
}
