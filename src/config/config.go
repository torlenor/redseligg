package config

// Config holds the complete AbyleBotter config
type Config struct {
	Bots struct {
		Discord struct {
			Enabled bool   `toml:"enabled"`
			Token   string `toml:"token"`
			Plugins struct {
				Echo struct {
					Enabled bool `toml:"enabled"`
				} `toml:"echo"`
				Giveaway struct {
					Enabled bool `toml:"enabled"`
				} `toml:"giveaway"`
			} `toml:"plugins"`
		} `toml:"discord"`
		Fake struct {
			Enabled bool `toml:"enabled"`
			Plugins struct {
				Echo struct {
					Enabled bool `toml:"enabled"`
				} `toml:"echo"`
				Giveaway struct {
					Enabled bool `toml:"enabled"`
				} `toml:"giveaway"`
			} `toml:"plugins"`
		} `toml:"fake"`
		Matrix struct {
			Enabled  bool   `toml:"enabled"`
			Server   string `toml:"server"`
			Username string `toml:"username"`
			Password string `toml:"password"`
			Token    string `toml:"token"`
			Plugins  struct {
				Echo struct {
					Enabled bool `toml:"enabled"`
				} `toml:"echo"`
				Giveaway struct {
					Enabled bool `toml:"enabled"`
				} `toml:"giveaway"`
			} `toml:"plugins"`
		} `toml:"matrix"`
	} `toml:"bots"`
}
