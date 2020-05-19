package slack

func (b *Bot) populateChannelList() error {
	conversations, err := b.getConversations()
	if err != nil {
		return err
	}

	for _, channel := range conversations {
		b.channels.addKnownChannel(channel)
	}

	b.log.Debugf("Added %d known channels", b.channels.Len())
	return nil
}

func (b *Bot) populateUserList() error {
	users, err := b.getUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		b.users.addKnownUser(user)
	}

	b.log.Debugf("Added %d known users", b.users.Len())
	return nil
}
