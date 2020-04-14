package versionplugin

import (
	"github.com/torlenor/abylebotter/plugin"
)

// VersionPlugin is a plugin implementing a version command
// which replies with the version of AbyleBotter
type VersionPlugin struct {
	plugin.AbyleBotterPlugin
}