package timedmessagesplugin

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/torlenor/redseligg/model"
	"github.com/torlenor/redseligg/storage"
	"github.com/torlenor/redseligg/storagemodels"
)

var now = time.Now

var errNotExist = errors.New("Timed message does not exist")

const (
	identFieldTimedMessages = "timedMessages"
)

func (p *TimedMessagesPlugin) getTimedMessages() (storagemodels.TimedMessagesPluginMessages, error) {
	s := p.getStorage()
	if s == nil {
		return storagemodels.TimedMessagesPluginMessages{}, ErrNoValidStorage
	}

	return s.GetTimedMessagesPluginMessages(p.BotID, p.PluginID, identFieldTimedMessages)
}

func (p *TimedMessagesPlugin) storeTimedMessages(data storagemodels.TimedMessagesPluginMessages) error {
	s := p.getStorage()
	if s == nil {
		p.API.LogError(ErrNoValidStorage.Error())
		return ErrNoValidStorage
	}

	err := s.StoreTimedMessagesPluginMessages(p.BotID, p.PluginID, identFieldTimedMessages, data)
	if err != nil {
		p.API.LogError(fmt.Sprintf("Error storing timed messages: %s", err))
		return fmt.Errorf("Error storing timed messages: %s", err)
	}

	return nil
}

func (p *TimedMessagesPlugin) addTimedMessage(channelID, message string, interval time.Duration) error {
	timedMessages, err := p.getTimedMessages()
	if err != nil && err != storage.ErrNotFound {
		return fmt.Errorf("Could not add new timed message: %s", err)
	}
	timedMessages.Messages = append(timedMessages.Messages, storagemodels.TimedMessagesPluginMessage{
		Text:      message,
		Interval:  interval,
		ChannelID: channelID,
	})

	p.API.LogTrace(fmt.Sprintf("Added message '%s' with interval %s for channel %s", message, interval, channelID))

	return p.storeTimedMessages(timedMessages)
}

func (p *TimedMessagesPlugin) removeTimedMessage(channelID, message string, interval time.Duration) error {
	timedMessages, err := p.getTimedMessages()
	if err != nil && err != storage.ErrNotFound {
		return fmt.Errorf("Could not remove timed message: %s", err)
	}

	var wasRemoved bool
	n := 0
	for _, x := range timedMessages.Messages {
		if !(x.ChannelID == channelID && x.Interval == interval && x.Text == message) {
			timedMessages.Messages[n] = x
			n++
		} else {
			p.API.LogTrace(fmt.Sprintf("Removed message '%s' with interval %s for channel %s", message, interval, channelID))
			wasRemoved = true
		}
	}
	timedMessages.Messages = timedMessages.Messages[:n]

	if !wasRemoved {
		return errNotExist
	}

	return p.storeTimedMessages(timedMessages)
}

func (p *TimedMessagesPlugin) removeAllTimedMessage(channelID, message string) error {
	timedMessages, err := p.getTimedMessages()
	if err != nil && err != storage.ErrNotFound {
		return fmt.Errorf("Could not remove timed message: %s", err)
	}

	var wasRemoved bool
	n := 0
	for _, x := range timedMessages.Messages {
		if !(x.ChannelID == channelID && x.Text == message) {
			timedMessages.Messages[n] = x
			n++
		} else {
			p.API.LogTrace(fmt.Sprintf("Removed message '%s' with interval %s for channel %s", message, x.Interval, channelID))
			wasRemoved = true
		}
	}
	timedMessages.Messages = timedMessages.Messages[:n]

	if !wasRemoved {
		return errNotExist
	}

	return p.storeTimedMessages(timedMessages)
}

func parseTimeStringToDuration(timeStr string) (time.Duration, error) {
	duration, err := time.ParseDuration(timeStr)
	if err != nil {
		return time.Duration(0), fmt.Errorf("Not a valid duration")
	}

	return duration, nil
}

func splitTmCommand(text string) (c string, interval time.Duration, msg string, err error) {
	var re = regexp.MustCompile(`(?m)^!tm +(add|remove) +([0-9]+[a-zA-Z]*) +(.+)$`)

	const cgCommand = 1
	const cgInterval = 2
	const cgMessage = 3

	matches := re.FindAllStringSubmatch(text, -1)

	if matches == nil || len(matches) < 1 {
		err = errors.New("Not a valid command")
		return
	}

	if len(matches[0]) > cgMessage {
		c = matches[0][cgCommand]
		interval, err = parseTimeStringToDuration(matches[0][cgInterval])
		if err != nil {
			return
		}
		msg = matches[0][cgMessage]
	} else {
		err = errors.New("Not a valid command")
	}

	return
}

// onCommand handles a !tm command.
func (p *TimedMessagesPlugin) onCommand(post model.Post) {
	if post.Content == "!tm add" {
		p.returnHelpAdd(post.ChannelID)
		return
	} else if post.Content == "!tm remove" {
		p.returnHelpRemove(post.ChannelID)
		return
	}

	if strings.HasPrefix(post.Content, "!tm remove all ") {
		cont := strings.Split(post.Content, " ")
		if len(cont) < 4 {
			p.returnHelpRemove(post.ChannelID)
			return
		}

		msg := strings.Join(cont[3:], " ")

		err := p.removeAllTimedMessage(post.ChannelID, msg)
		if err == errNotExist {
			p.returnMessage(post.ChannelID, "Timed message to remove does not exist.")
			return
		} else if err != nil {
			p.API.LogError(fmt.Sprintf("Could not remove all timed messages: %s", err))
			p.returnMessage(post.ChannelID, fmt.Sprintf("Could not remove all timed message. Please try again later."))
			return
		}
		p.returnMessage(post.ChannelID, fmt.Sprintf("All timed messages with text '%s' removed.", msg))
		return
	}

	c, interval, message, err := splitTmCommand(post.Content)
	if err != nil {
		p.API.LogError(fmt.Sprintf("Error parsing !tm command '%s': %s", post.Content, err))
		p.returnHelp(post.ChannelID)
		return
	}

	switch c {
	case "add":
		err = p.addTimedMessage(post.ChannelID, message, interval)
	case "remove":
		err = p.removeTimedMessage(post.ChannelID, message, interval)
	}

	if err == errNotExist {
		p.returnMessage(post.ChannelID, "Timed message to remove does not exist.")
		return
	} else if err != nil {
		p.API.LogError(fmt.Sprintf("Could not %s timed message: %s", c, err))
		p.returnMessage(post.ChannelID, fmt.Sprintf("Could not %s timed message. Please try again later.", c))
		return
	}

	switch c {
	case "add":
		p.returnMessage(post.ChannelID, fmt.Sprintf("Timed message '%s' with interval %s added.", message, interval))
	case "remove":
		p.returnMessage(post.ChannelID, fmt.Sprintf("Timed message '%s' with interval %s removed.", message, interval))
	}
}
