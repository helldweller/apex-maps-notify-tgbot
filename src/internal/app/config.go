package app

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config is the main configuration structure of this application
type Config struct {
	ApexAPIKey     string        `env:"APEX_API_KEY"    env-required:"true"`
	Loglevel       string        `env:"LOG_LEVEL"       env-default:"error"`
	BotDebug       bool          `env:"TGBOT_DEBUG"     env-default:"false"`
	BotAPIKey      string        `env:"TGBOT_API_KEY"   env-required:"true"`
	HTTPProxy      string        `env:"HTTP_PROXY"      env-default:""`
	UpdateInterval time.Duration `env:"UPDATE_INTERVAL" env-default:"120s"`
}

// loadConfig reads the configuration from environment variables.
func loadConfig() (Config, error) {
	var c Config
	if err := cleanenv.ReadEnv(&c); err != nil {
		return Config{}, fmt.Errorf("reading configuration: %w", err)
	}
	return c, nil
}
