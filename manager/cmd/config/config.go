package config

// Config is the structure for the configuration data read from the env.
type Config struct {
	DataSource string `envconfig:"DATA_SOURCE" default:"app.db"`

	// Logging.
	LogLevel             int8 `envconfig:"LOG_LEVEL" default:"0"` // trace: -1, debug: 0, info: 1, warn: 2
	HumanFriendlyLogging bool `envconfig:"HUMAN_FRIENDLY_LOGGING" default:"false"`
}

// NewConfig returns a config instance initialized with default values.
func NewConfig() *Config {
	return &Config{}
}
