package httppingplugin

import (
	"github.com/torlenor/redseligg/plugin"
)

// HTTPPingPlugin is a plugin implementing a httpping command
// which sends back the results of a http connection atttempt
// to the url provided with the command.
type HTTPPingPlugin struct {
	plugin.AbyleBotterPlugin
}
