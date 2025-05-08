package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes the logger
func Init() {
	// Configure zerolog
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// Pretty logging for development
	if os.Getenv("APP_ENV") != "prod" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	}

	// Add caller info
	log.Logger = log.With().Caller().Logger()
}

// Get returns the global logger instance
func Get() *zerolog.Logger {
	return &log.Logger
}
