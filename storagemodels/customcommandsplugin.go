// Package storagemodels contains custom storage models used by plugins
package storagemodels

// CustomCommandsPluginCommand is a CustomCommandsPlugin storage model.
// It represents one custom command.
type CustomCommandsPluginCommand struct {
	Command string
	Text    string

	ChannelID string
}

// CustomCommandsPluginCommands is a CustomCommandsPlugin storage model.
// It is used to store all custom commands.
type CustomCommandsPluginCommands struct {
	Commands []CustomCommandsPluginCommand
}
