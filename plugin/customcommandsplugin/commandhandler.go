package customcommandsplugin

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/torlenor/abylebotter/model"
	"github.com/torlenor/abylebotter/storage"
	"github.com/torlenor/abylebotter/storagemodels"
)

var now = time.Now

var errNotExist = errors.New("Timed message does not exist")

const (
	identField = "customCommands"
)

func (p *CustomCommandsPlugin) getCommands() (storagemodels.CustomCommandsPluginCommands, error) {
	s := p.getStorage()
	if s == nil {
		return storagemodels.CustomCommandsPluginCommands{}, ErrNoValidStorage
	}

	return s.GetCustomCommandsPluginCommands(p.BotID, p.PluginID, identField)
}

func (p *CustomCommandsPlugin) storeCommands(data storagemodels.CustomCommandsPluginCommands) error {
	s := p.getStorage()
	if s == nil {
		p.API.LogError(ErrNoValidStorage.Error())
		return ErrNoValidStorage
	}

	err := s.StoreCustomCommandsPluginCommands(p.BotID, p.PluginID, identField, data)
	if err != nil {
		p.API.LogError(fmt.Sprintf("Error storing timed messages: %s", err))
		return fmt.Errorf("Error storing timed messages: %s", err)
	}

	return nil
}

func (p *CustomCommandsPlugin) addCommand(channelID, customCommand string, message string) error {
	commands, err := p.getCommands()
	if err != nil && err != storage.ErrNotFound {
		return fmt.Errorf("Could not add new custom command: %s", err)
	}

	updated := false
	for i, c := range commands.Commands {
		if c.Command == customCommand {
			updated = true
			commands.Commands[i].Text = message
		}
	}

	if !updated {
		commands.Commands = append(commands.Commands, storagemodels.CustomCommandsPluginCommand{
			Text:      message,
			Command:   customCommand,
			ChannelID: channelID,
		})
		p.API.LogTrace(fmt.Sprintf("Added custom command '%s' with message '%s' for channel %s", customCommand, message, channelID))
	} else {
		p.API.LogTrace(fmt.Sprintf("Updated custom command '%s' with message '%s' for channel %s", customCommand, message, channelID))
	}
	return p.storeCommands(commands)
}

func (p *CustomCommandsPlugin) removeCommand(channelID, customCommand string) error {
	commands, err := p.getCommands()
	if err != nil && err != storage.ErrNotFound {
		return fmt.Errorf("Could not remove custom command: %s", err)
	}

	var wasRemoved bool
	n := 0
	for _, x := range commands.Commands {
		if !(x.ChannelID == channelID && x.Command == customCommand) {
			commands.Commands[n] = x
			n++
		} else {
			p.API.LogTrace(fmt.Sprintf("Removed command '%s' for channel %s", customCommand, channelID))
			wasRemoved = true
		}
	}
	commands.Commands = commands.Commands[:n]

	if !wasRemoved {
		return errNotExist
	}

	return p.storeCommands(commands)
}

func splitCommand(text string) (c string, customCommand string, msg string, err error) {
	var re = regexp.MustCompile(`(?m)^!customcommand +(add|remove) +([a-zA-Z]+) *(.*)$`)

	const cgCommand = 1
	const cgCustomCommand = 2
	const cgMessage = 3

	matches := re.FindAllStringSubmatch(strings.Trim(text, " "), -1)

	if matches == nil || len(matches) < 1 {
		err = errors.New("Not a valid command")
		return
	}
	if len(matches[0]) > cgCustomCommand {
		switch matches[0][cgCommand] {
		case "add":
			if len(matches[0]) > cgMessage && len(matches[0][cgMessage]) > 0 {
				msg = matches[0][cgMessage]
			} else {
				err = errors.New("Not a valid command: No message provided")
				return
			}
		case "remove":
			// do nothing
		default:
			err = errors.New("Not a valid command. Wrong command: Only add and remove are supported")
			return
		}
		c = matches[0][cgCommand]
		customCommand = matches[0][cgCustomCommand]
	} else {
		err = errors.New("Not a valid command")
	}

	return
}

// onCommand handles a !customcommand command.
func (p *CustomCommandsPlugin) onCommand(post model.Post) {
	if post.Content == "!customcommand add" {
		p.returnHelpAdd(post.ChannelID)
		return
	} else if post.Content == "!customcommand remove" {
		p.returnHelpRemove(post.ChannelID)
		return
	}

	c, customCommand, message, err := splitCommand(post.Content)
	if err != nil {
		p.API.LogError(fmt.Sprintf("Error parsing !customcommand command '%s': %s", post.Content, err))
		p.returnHelp(post.ChannelID)
		return
	}

	switch c {
	case "add":
		err = p.addCommand(post.ChannelID, customCommand, message)
	case "remove":
		err = p.removeCommand(post.ChannelID, customCommand)
	}

	if err == errNotExist {
		p.returnMessage(post.ChannelID, "Custom command to remove does not exist.")
		return
	} else if err != nil {
		p.API.LogError(fmt.Sprintf("Could not %s custom command: %s", c, err))
		p.returnMessage(post.ChannelID, fmt.Sprintf("Could not %s custom command. Please try again later.", c))
		return
	}

	switch c {
	case "add":
		p.returnMessage(post.ChannelID, fmt.Sprintf("Custom command '%s' with message '%s' added.", customCommand, message))
	case "remove":
		p.returnMessage(post.ChannelID, fmt.Sprintf("Custom command '%s' removed.", customCommand))
	}
}
