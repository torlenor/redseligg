package discord

func combineUsernameAndDiscriminator(username, discriminator string) string {
	return username + "#" + discriminator
}
