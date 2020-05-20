package commanddispatcher

import (
	"strings"

	"github.com/torlenor/redseligg/logging"
	"github.com/torlenor/redseligg/model"
)

var log = logging.Get("CommandDispatcher")

var defaultCallPrefix = "!"

type receiver interface {
	// OnCommand delivers the command, the whitespace-trimmed content where the command is already stripped off and the raw Post.
	OnCommand(cmd string, content string, post model.Post)
}

// CommandDispatcher provides an architecture to let plugins (or other entities) register commands and get notified.
// At first we will only support one receiver for a specific command.
type CommandDispatcher struct {
	callPrefix string

	receivers map[string]receiver // [cmd]
}

// New CommandDispatcher
func New(callPrefix string) *CommandDispatcher {
	var c CommandDispatcher
	if len(callPrefix) == 0 {
		c.callPrefix = defaultCallPrefix
	} else {
		c.callPrefix = callPrefix
	}

	log.Tracef("Created new CommandDispatcher with call prefix = '%s'", c.callPrefix)

	c.receivers = make(map[string]receiver)

	return &c
}

// Register a new command receiver with the specified command (without call prefix).
func (c *CommandDispatcher) Register(cmd string, r receiver) {
	log.Tracef("Registering command %s", cmd)
	if len(cmd) > 0 {
		c.receivers[cmd] = r
	} else {
		log.Warn("Tried to register an empty command")
	}
}

// Unregister removes the specified command from the receivers list if it exists.
func (c *CommandDispatcher) Unregister(cmd string) {
	delete(c.receivers, cmd)
}

// OnPost feeds a post to the CommandDispatcher which will then do its magic
func (c *CommandDispatcher) OnPost(post model.Post) {
	if len(post.Content) < 2 {
		return
	}
	splitted := strings.Split(post.Content, " ")
	if !strings.HasPrefix(splitted[0], c.callPrefix) {
		return
	}
	cmd := splitted[0][1:]
	content := ""
	if len(splitted) > 1 {
		content = strings.Join(splitted[1:], " ")
		content = strings.TrimSpace(content)
	}
	for c, r := range c.receivers {
		if cmd == c {
			r.OnCommand(cmd, content, post)
		}
	}
}
