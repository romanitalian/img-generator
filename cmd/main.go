package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/romanitalian/img-generate/v2/configs"
	"github.com/romanitalian/img-generate/v2/internal/server"
	"github.com/romanitalian/img-generate/v2/pkg/logger"
)

var confPath = flag.String("conf-path", ".env", "Path to config env.")

func main() {
	flag.Parse()

	// Initialize logger
	logger.Init()
	log := logger.Get()

	// Create base context for the application
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	go func() {
		sig := <-sigChan
		log.Info().Str("signal", sig.String()).Msg("received shutdown signal")
		cancel()
	}()

	log.Info().Msg("starting image generator service")

	conf, err := configs.New(*confPath)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	if err := server.Run(ctx, conf); err != nil {
		log.Fatal().Err(err).Msg("server error")
	}
}
