package rollplugin

import (
	"time"

	"github.com/torlenor/abylebotter/plugin"
)

type randomizer interface {
	random(max int) int
}

// RollPlugin is a plugin implementing a roll command
// which sends back a random number in the range [0,100] per default or
// [0,specific_number] if the roll command contained a number.
type RollPlugin struct {
	plugin.AbyleBotterPlugin
	randomizer randomizer
}

// New returns a new RollPlugin
func New() (RollPlugin, error) {
	ep := RollPlugin{
		randomizer: newRoller(time.Now().UnixNano()),
	}

	return ep, nil
}
