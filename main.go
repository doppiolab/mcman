package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/doppiolab/mcman/internal/config"
	"github.com/doppiolab/mcman/internal/server"
	"github.com/rs/zerolog/log"
)

var (
	configFileName = flag.String("config", "default.yaml", "config file name")
)

func main() {
	flag.Parse()
	log.Info().Str("config-file", *configFileName).Msg("start mcman")

	cfg := config.MustGetConfig(*configFileName)

	// launch server
	svr, err := server.New(&cfg.Server)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create server")
	}
	go (func() {
		log.Info().Msgf("start to listen on %s", svr.Addr)
		if err := svr.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error().Err(err).Msg("cannot serve http server")
		}
	})()

	// wait for shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)
	<-quit

	// graceful shutdown
	log.Info().Msg("stopping http server")
	ctx := context.Background()
	if err := svr.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("cannot shutdown gateway server")
	}
}
