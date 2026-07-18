package app

import (
	"os"

	"github.com/sirupsen/logrus"
)

// log is usable before configuration (default logrus settings) so that any
// code path — including tests — can log without an explicit setup step.
var log = logrus.New()

// configureLogger applies the JSON formatter, stdout output and the configured
// log level to the global logger.
func configureLogger(cfg Config) {
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)

	level, err := logrus.ParseLevel(cfg.Loglevel)
	if err != nil {
		log.SetLevel(logrus.ErrorLevel)
		log.Error("Can`t parse LOG_LEVEL. Used default value: LOG_LEVEL=error")
	} else {
		log.SetLevel(level)
	}

	log.Infof("Used json logger and loglevel: %s", cfg.Loglevel)
}
