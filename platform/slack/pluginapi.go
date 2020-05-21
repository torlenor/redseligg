package slack

import (
	"fmt"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/utils"
)

var version string

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
func (b *Bot) CreatePost(post model.Post) (model.PostResponse, error) {
	var response genericChatResponse

	if post.IsPrivate {
		var userID string
		if len(post.User.ID) != 0 {
			userID = post.User.ID
		} else if len(post.User.Name) != 0 {
			var err error
			userID, err = b.users.getUserNameByID(post.User.Name)
			if err != nil {
				return model.PostResponse{}, fmt.Errorf("User not found, not sending Whisper")
			}
		} else {
			return model.PostResponse{}, fmt.Errorf("Plugin did not provide User or UserID, not sending Whisper")
		}
		var err error
		response, err = b.sendWhisper(userID, post.Content)
		if err != nil {
			return model.PostResponse{}, fmt.Errorf("Error sending whisper: %s", err)
		}
	} else {
		var err error
		response, err = b.sendMessage(post.ChannelID, post.Content)
		if err != nil {
			return model.PostResponse{}, fmt.Errorf("Error sending message: %s", err)
		}
	}

	b.log.Debugf("Got response from message/whisper sending: %v", response)

	return model.PostResponse{PostedMessageIdent: model.MessageIdentifier{ID: response.TS, Channel: response.Channel}}, nil
}

// UpdatePost updates a post.
func (b *Bot) UpdatePost(messageID model.MessageIdentifier, newPost model.Post) (model.PostResponse, error) {
	response, err := b.chatUpdate(messageID.Channel, messageID.ID, newPost.Content)
	if err != nil {
		return model.PostResponse{}, fmt.Errorf("Error updating post: %s", err)
	}

	return model.PostResponse{PostedMessageIdent: model.MessageIdentifier{ID: response.TS, Channel: response.Channel}}, nil
}

// DeletePost deletes a post.
func (b *Bot) DeletePost(messageID model.MessageIdentifier) (model.PostResponse, error) {
	response, err := b.chatDelete(messageID.Channel, messageID.ID)
	if err != nil {
		return model.PostResponse{}, fmt.Errorf("Error deleting post: %s", err)
	}

	return model.PostResponse{PostedMessageIdent: model.MessageIdentifier{ID: response.TS, Channel: response.Channel}}, nil
}

// GetReaction gives back the platform specific string for a reaction, e.g., one -> :one:
func (b *Bot) GetReaction(reactionName string) (string, error) {
	return "", fmt.Errorf("Not implemented")
}

// LogTrace writes a log message to the server log file.
func (b *Bot) LogTrace(msg string) {
	b.log.Tracef("Error from plugin: %s", msg)
}

// LogDebug writes a log message to the server log file.
func (b *Bot) LogDebug(msg string) {
	b.log.Debugf("From plugin: %s", msg)
}

// LogInfo writes a log message to the server log file.
func (b *Bot) LogInfo(msg string) {
	b.log.Infof("From plugin: %s", msg)
}

// LogWarn writes a log message to the server log file.
func (b *Bot) LogWarn(msg string) {
	b.log.Warnf("From plugin: %s", msg)
}

// LogError writes a log message to the server log file.
func (b *Bot) LogError(msg string) {
	b.log.Errorf("From plugin: %s", msg)
}

// GetVersion returns the version of the server.
func (b *Bot) GetVersion() string {
	return utils.Version().Get() + " (" + utils.Version().GetCompTime() + ")"
}
