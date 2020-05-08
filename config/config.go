package config

// API holds the API settings for the Redseligg configuration
type API struct {
	Enabled bool `toml:"enabled"`
	// IP Address the REST API listens on
	// If empty or non-existing listen on all interfaces
	IP string `toml:"ip"`
	// Port the REST API listens on
	Port string `toml:"port"`
}
