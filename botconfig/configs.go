package botconfig

// DiscordConfig contains config related to the Discord component
type DiscordConfig struct {
	ID     string `toml:"id" json:"id"`
	Token  string `toml:"token" json:"token"`
	Secret string `toml:"secret" json:"secret"`
}

// MatrixConfig contains config related to the Matrix component
type MatrixConfig struct {
	Server   string `toml:"server" json:"server"`
	Username string `toml:"username" json:"username"`
	Password string `toml:"password" json:"password"`
}

// MattermostConfig contains config related to the Mattermost component
type MattermostConfig struct {
	Server   string `toml:"server" json:"server"`
	Username string `toml:"username" json:"username"`
	Password string `toml:"password" json:"password"`
}

// SlackConfig contains config related to the Mattermost component
type SlackConfig struct {
	Workspace string `toml:"workspace" json:"workspace"`
	Token     string `toml:"token" json:"token"`
}

// TwitchConfig contains config related to the Twitch component
type TwitchConfig struct {
	Username string   `toml:"username" json:"username"`
	Token    string   `toml:"token" json:"token"`
	Channels []string `toml:"channels" json:"channels"`
}
