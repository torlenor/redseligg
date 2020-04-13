package plugin

import (
	"fmt"

	"github.com/torlenor/abylebotter/model"
)

// The MockAPI can be used for testing Plugins by providing helper functions.
// It mimics all the functions a real bot would, but in addition provides helper functions
// for unit tests.
type MockAPI struct {
	WasCreatePostCalled bool
	LastCreatePostPost  model.Post

	WasUpdatePostCalled     bool
	LastUpdatePostMessageID model.MessageIdentifier
	LastUpdatePostPost      model.Post

	PostResponse model.PostResponse

	ErrorToReturn error
}

// Reset the MockAPI
func (b *MockAPI) Reset() {
	b.WasCreatePostCalled = false
	b.LastCreatePostPost = model.Post{}

	b.WasUpdatePostCalled = false
	b.LastUpdatePostMessageID = model.MessageIdentifier{}
	b.LastUpdatePostPost = model.Post{}
}

// RegisterCommand registers a custom slash "/" or "!" command, depending on what the bot supports.
func (b *MockAPI) RegisterCommand(command string) error { return nil }

// UnregisterCommand unregisters a command previously registered via RegisterCommand.
func (b *MockAPI) UnregisterCommand(command string) error { return nil }

// GetUsers a list of all users the bot knows.
func (b *MockAPI) GetUsers() ([]model.User, error) { return nil, nil }

// GetUser gets a user by their ID.
func (b *MockAPI) GetUser(userID string) (model.User, error) { return model.User{}, nil }

// GetUserByUsername gets a user by their name.
func (b *MockAPI) GetUserByUsername(name string) (model.User, error) { return model.User{}, nil }

// GetChannel gets a channel by its ID.
func (b *MockAPI) GetChannel(channelID string) (model.Channel, error) { return model.Channel{}, nil }

// GetChannelByName gets a channel by its name.
func (b *MockAPI) GetChannelByName(name string) (model.Channel, error) { return model.Channel{}, nil }

// CreatePost creates a post.
func (b *MockAPI) CreatePost(post model.Post) (model.PostResponse, error) {
	b.WasCreatePostCalled = true
	b.LastCreatePostPost = post
	return b.PostResponse, b.ErrorToReturn
}

// UpdatePost updates a post.
func (b *MockAPI) UpdatePost(messageID model.MessageIdentifier, newPost model.Post) (model.PostResponse, error) {
	b.WasUpdatePostCalled = true
	b.LastUpdatePostPost = newPost
	b.LastUpdatePostMessageID = messageID
	return b.PostResponse, b.ErrorToReturn
}

// DeletePost deletes a post.
func (b *MockAPI) DeletePost(messageID model.MessageIdentifier) (model.PostResponse, error) {
	return model.PostResponse{}, fmt.Errorf("Not implemented")
}

// GetReaction gives back the platform specific string for a reaction, e.g., one -> :one:
func (b *MockAPI) GetReaction(reactionName string) (string, error) {
	return "", fmt.Errorf("Not implemented")
}

// LogTrace writes a log message to the server log file.
func (b *MockAPI) LogTrace(msg string) {}

// LogDebug writes a log message to the server log file.
func (b *MockAPI) LogDebug(msg string) {}

// LogInfo writes a log message to the server log file.
func (b *MockAPI) LogInfo(msg string) {}

// LogWarn writes a log message to the server log file.
func (b *MockAPI) LogWarn(msg string) {}

// LogError writes a log message to the server log file.
func (b *MockAPI) LogError(msg string) {}

// GetVersion returns the version of the server.
func (b *MockAPI) GetVersion() string {
	return "0.0.1"
}
