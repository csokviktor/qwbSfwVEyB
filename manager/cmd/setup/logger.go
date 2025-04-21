package setup

import (
	"os"

	"github.com/csokviktor/lib_manager/cmd/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func Logger(cfg *config.Config) {
	zerolog.TimestampFieldName = "datetime"
	zerolog.LevelFieldName = "loglevel"
	zerolog.MessageFieldName = "message"
	zerolog.SetGlobalLevel(zerolog.Level(cfg.LogLevel))
	log.Logger = log.With().Str("component", "lib_manager").Logger()
	if cfg.HumanFriendlyLogging {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
