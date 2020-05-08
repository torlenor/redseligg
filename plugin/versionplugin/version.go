package versionplugin

import (
	"github.com/torlenor/redseligg/plugin"
)

// VersionPlugin is a plugin implementing a version command
// which replies with the version of Redseligg
type VersionPlugin struct {
	plugin.RedseliggPlugin
}
