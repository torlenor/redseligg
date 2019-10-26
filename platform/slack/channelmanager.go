package slack

import "fmt"

type channelManager struct {
	knownChannels     map[string]channel // key is ChannelID
	knownChannelNames map[string]string  // mapping of ChannelName to ChannelID
	knownChannelIDs   map[string]string  // mapping of ChannelID to UserChannelNameName
}

func newChannelManager() channelManager {
	return channelManager{
		knownChannels:     make(map[string]channel),
		knownChannelNames: make(map[string]string),
		knownChannelIDs:   make(map[string]string),
	}
}

func (cm *channelManager) addKnownChannel(channel channel) {
	cm.knownChannels[channel.ID] = channel
	cm.knownChannelNames[channel.Name] = channel.ID
	cm.knownChannelIDs[channel.ID] = channel.Name
}

func (cm channelManager) getChannelByID(id string) (channel, error) {
	if channel, ok := cm.knownChannels[id]; ok {
		return channel, nil
	}
	return channel{}, fmt.Errorf("Channel with ID %s not known", id)
}

func (cm channelManager) getChannelByName(name string) (channel, error) {
	if id, ok := cm.knownChannelNames[name]; ok {
		return cm.knownChannels[id], nil
	}
	return channel{}, fmt.Errorf("Channel with Name %s not known", name)
}

func (cm channelManager) getChannelNameByID(id string) (string, error) {
	if name, ok := cm.knownChannelIDs[id]; ok {
		return name, nil
	}
	return "", fmt.Errorf("Channel with ID %s not known", id)
}

func (cm channelManager) getChannelIDByName(name string) (string, error) {
	if id, ok := cm.knownChannelNames[name]; ok {
		return id, nil
	}
	return "", fmt.Errorf("Channel with Name %s not known", name)
}

func (cm channelManager) isChannelIDKnown(id string) bool {
	_, ok := cm.knownChannelIDs[id]
	return ok
}

func (cm channelManager) isChannelNameKnown(name string) bool {
	_, ok := cm.knownChannelNames[name]
	return ok
}

func (cm channelManager) Len() int {
	return len(cm.knownChannels)
}
