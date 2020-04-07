package slack

import (
	"fmt"

	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/utils"
)

var version string

// RegisterCommand registers a custom slash or ! command, depending on what the bot supports.
func (b *Bot) RegisterCommand(command string) error { return nil }

// UnregisterCommand unregisters a command previously registered via RegisterCommand.
func (b *Bot) UnregisterCommand(command string) error { return nil }

// GetUsers a list of users based on search options.
func (b *Bot) GetUsers() ([]model.User, error) { return nil, nil }

// GetUser gets a user.
func (b *Bot) GetUser(userID string) (model.User, error) { return model.User{}, nil }

// GetUserByUsername gets a user by their username.
func (b *Bot) GetUserByUsername(name string) (model.User, error) { return model.User{}, nil }

// GetChannel gets a channel.
func (b *Bot) GetChannel(channelID string) (model.Channel, error) { return model.Channel{}, nil }

// GetChannelByName gets a channel by its name, given a team id.
func (b *Bot) GetChannelByName(name string) (model.Channel, error) { return model.Channel{}, nil }

// CreatePost creates a post.
func (b *Bot) CreatePost(post model.Post) error {

	if post.IsPrivate {
		var userID string
		if len(post.UserID) != 0 {
			userID = post.UserID
		} else if len(post.User) != 0 {
			var err error
			userID, err = b.users.getUserNameByID(post.User)
			if err != nil {
				return fmt.Errorf("User not found, not sending Whisper")
			}
		} else {
			return fmt.Errorf("Plugin did not provide User or UserID, not sending Whisper")
		}
		err := b.sendWhisper(userID, post.Content)
		if err != nil {
			return fmt.Errorf("Error sending whisper: %s", err)
		}
	} else {
		err := b.sendMessage(post.ChannelID, post.Content)
		if err != nil {
			return fmt.Errorf("Error sending message: %s", err)
		}
	}

	return nil
}

// LogTrace writes a log message to the server log file.
func (b *Bot) LogTrace(msg string) {}

// LogDebug writes a log message to the server log file.
func (b *Bot) LogDebug(msg string) {}

// LogInfo writes a log message to the server log file.
func (b *Bot) LogInfo(msg string) {}

// LogWarn writes a log message to the server log file.
func (b *Bot) LogWarn(msg string) {}

// LogError writes a log message to the server log file.
func (b *Bot) LogError(msg string) {}

// GetVersion returns the version of the server.
func (b *Bot) GetVersion() string {
	return utils.Version().Get() + " (" + utils.Version().GetCompTime() + ")"
}
