package app

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ApexApiKey string `env:"APEX_API_KEY"    env-required:"true"`
	Loglevel   string `env:"LOG_LEVEL"       env-default:"error"`
	BotDebug   bool   `env:"TGBOT_DEBUG"     env-default:"false"`
	BotApiKey  string `env:"TGBOT_API_KEY"   env-required:"true"`
}

var conf Config

func init() {
	err := cleanenv.ReadEnv(&conf)
	if err != nil {
		fmt.Printf("Something went wrong while reading the configuration: %s", err)
		os.Exit(1)
	}
}
