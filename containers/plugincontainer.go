package containers

import (
	"sync"

	"github.com/torlenor/abylebotter/events"
	"github.com/torlenor/abylebotter/plugins"
)

// PluginContainer is used to manage a collection of plugins
// and send/receive messages to/from them.
type PluginContainer struct {
	mutex sync.RWMutex

	receivers    map[plugins.Plugin]chan events.ReceiveMessage
	knownPlugins []plugins.Plugin

	sendMessageChan chan events.SendMessage
}

// Add a new plugin to the container
func (p *PluginContainer) Add(plugin plugins.Plugin) {
	plugin.ConnectChannels(p.receiveChannel(plugin), p.SendChannel())
	p.knownPlugins = append(p.knownPlugins, plugin)
}

// RemoveAll removes all plugins from the container essentially disconnecting and
// throwing away all plugins
func (p *PluginContainer) RemoveAll() {
	p.mutex.Lock()
	for _, pluginChannel := range p.receivers {
		close(pluginChannel)
	}
	p.receivers = nil
	p.knownPlugins = nil
	p.mutex.Unlock()
}

// Size returns the number of plugin currently in the container
func (p *PluginContainer) Size() int {
	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return len(p.knownPlugins)
}

// SendChannel gives you the channel where Plugin messages are received
func (p *PluginContainer) SendChannel() chan events.SendMessage {
	if p.sendMessageChan == nil {
		p.mutex.Lock()
		p.sendMessageChan = make(chan events.SendMessage)
		p.mutex.Unlock()
	}

	p.mutex.RLock()
	defer p.mutex.RUnlock()
	return p.sendMessageChan
}

// Send a message to all connected plugins (blocking)
func (p *PluginContainer) Send(receiveMessage events.ReceiveMessage) {
	p.mutex.RLock()
	for _, pluginChannel := range p.receivers {
		select {
		case pluginChannel <- receiveMessage:
		default:
		}
	}
	p.mutex.RUnlock()
}

// receiveChannel returns the channel which is used to notify
// about received messages from the bot
func (p *PluginContainer) receiveChannel(plugin plugins.Plugin) <-chan events.ReceiveMessage {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.receivers == nil {
		p.receivers = make(map[plugins.Plugin]chan events.ReceiveMessage)
	}

	p.receivers[plugin] = make(chan events.ReceiveMessage)
	return p.receivers[plugin]
}
