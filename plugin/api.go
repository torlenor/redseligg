package plugin

import (
	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/storage"
)

// The API can be used to retrieve data or perform actions on behalf of the plugin.
//
// A valid Bot has to implement all of these functions (or provide a wrapper for the Plugin).
//
// Plugins obtain access to this API by embedding RedseliggPlugin.
type API interface {
	// HasFeature returns true if the bot serving the API implements the feature
	HasFeature(feature string) bool

	// GetStorage returns the storage or nil if none is provided by the platform
	GetStorage() storage.Storage

	// RegisterCommand registers a custom slash "/" or "!" command, depending on what the bot supports.
	RegisterCommand(p Hooks, command string) error

	// UnRegisterCommand unregisters a command previously registered via RegisterCommand.
	UnRegisterCommand(command string) error

	// GetUsers a list of all users the bot knows.
	GetUsers() ([]model.User, error)

	// GetUser gets a user by their ID.
	GetUser(userID string) (model.User, error)

	// GetUserByUsername gets a user by their name.
	GetUserByUsername(name string) (model.User, error)

	// GetChannel gets a channel by its ID.
	GetChannel(channelID string) (model.Channel, error)

	// GetChannelByName gets a channel by its name.
	GetChannelByName(name string) (model.Channel, error)

	// CreatePost creates a post.
	CreatePost(post model.Post) (model.PostResponse, error)

	// UpdatePost updates a previous post.
	// The messageID must be a valid model.MessageIdentifier.
	// Currently such a messageID is supplied in CreatePost calls when the platform supports it.
	UpdatePost(messageID model.MessageIdentifier, newPost model.Post) (model.PostResponse, error)

	// DeletePost deletes a previous post.
	// The messageID must be a valid model.MessageIdentifier.
	// Currently such a messageID is supplied in CreatePost calls when the platform supports it.
	DeletePost(messageID model.MessageIdentifier) (model.PostResponse, error)

	// GetReaction gives back the platform specific string for a reaction, e.g., one -> :one:
	// If the reaction is unknown returns an error.
	GetReaction(reactionName string) (string, error)

	// LogTrace writes a log message to the server log file.
	LogTrace(msg string)

	// LogDebug writes a log message to the server log file.
	LogDebug(msg string)

	// LogInfo writes a log message to the server log file.
	LogInfo(msg string)

	// LogWarn writes a log message to the server log file.
	LogWarn(msg string)

	// LogError writes a log message to the server log file.
	LogError(msg string)

	// GetVersion returns the version of the server.
	GetVersion() string
}
