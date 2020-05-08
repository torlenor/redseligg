package quotesplugin

import (
	"fmt"

	"github.com/torlenor/redseligg/botconfig"
)

type config struct {
	Mods     []string
	OnlyMods bool
}

func parseConfig(c botconfig.PluginConfig) (config, error) {
	if c.Type != PLUGIN_TYPE {
		return config{}, fmt.Errorf("Not a " + PLUGIN_TYPE + " plugin config")
	}

	var mods []string
	var onlyMods bool

	var ok bool

	if modEntries, ok := c.Config["mods"].([]interface{}); ok {
		for _, modEntry := range modEntries {
			if mod, ok := modEntry.(string); ok {
				mods = append(mods, mod)
			}
		}
	} else if modEntries, ok := c.Config["mods"].([]string); ok {
		mods = modEntries
	}

	if onlyMods, ok = c.Config["onlymods"].(bool); !ok {
	}

	if onlyMods && len(mods) == 0 {
		return config{}, fmt.Errorf("Cannot have a " + PLUGIN_TYPE + "  plugin configuration with OnlyMods = true but no Mods defined")
	}

	cfg := config{
		Mods:     mods,
		OnlyMods: onlyMods,
	}

	return cfg, nil
}
