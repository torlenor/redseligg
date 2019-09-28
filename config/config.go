package config

// API holds the API settings for the AbyleBotter configuration
type API struct {
	Enabled bool `toml:"enabled"`
	// IP Address the REST API listens on
	// If empty or non-existing listen on all interfaces
	IP string `toml:"ip"`
	// Port the REST API listens on
	Port string `toml:"port"`
}

// General holds the general settings for the AbyleBotter configuration
type General struct {
	API API `toml:"api"`
}

// Plugins holds the plugins part of the AbyleBotter configuration
type Plugins struct {
	Echo struct {
		Enabled      bool `toml:"enabled"`
		OnlyWhispers bool `toml:"onlywhispers"`
	} `toml:"echo"`
	Giveaway struct {
		Enabled bool `toml:"enabled"`
	} `toml:"giveaway"`
	SendMessage struct {
		Enabled bool `toml:"enabled"`
	} `toml:"sendmessages"`
}

// MattermostConfig contains config related to the Mattermost component
type MattermostConfig struct {
	Enabled  bool    `toml:"enabled"`
	Server   string  `toml:"server"`
	Username string  `toml:"username"`
	Password string  `toml:"password"`
	UseToken bool    `toml:"usetoken"`
	Token    string  `toml:"token"`
	Plugins  Plugins `toml:"plugins"`
}

// SlackConfig contains config related to the Mattermost component
type SlackConfig struct {
	Enabled   bool    `toml:"enabled"`
	Workspace string  `toml:"workspace"`
	Token     string  `toml:"token"`
	Plugins   Plugins `toml:"plugins"`
}

// Config holds the complete AbyleBotter config
type Config struct {
	General General `toml:"general"`
	Bots    struct {
		Discord struct {
			Enabled bool    `toml:"enabled"`
			ID      string  `toml:"id"`
			Token   string  `toml:"token"`
			Secret  string  `toml:"secret"`
			Plugins Plugins `toml:"plugins"`
		} `toml:"discord"`
		Fake struct {
			Enabled bool    `toml:"enabled"`
			Plugins Plugins `toml:"plugins"`
		} `toml:"fake"`
		Matrix struct {
			Enabled  bool    `toml:"enabled"`
			Server   string  `toml:"server"`
			Username string  `toml:"username"`
			Password string  `toml:"password"`
			Token    string  `toml:"token"`
			Plugins  Plugins `toml:"plugins"`
		} `toml:"matrix"`
		Mattermost MattermostConfig `toml:"mattermost"`
		Slack      SlackConfig      `toml:"slack"`
	} `toml:"bots"`
}
