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
	"github.com/doppiolab/mcman/internal/minecraft"
	"github.com/doppiolab/mcman/internal/minecraft/logstream"
	"github.com/doppiolab/mcman/internal/minecraft/logstream/callback"
	"github.com/doppiolab/mcman/internal/minecraft/world"
	"github.com/doppiolab/mcman/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	configFileName = flag.String("config", "default.yaml", "config file name")
)

func main() {
	flag.Parse()

	cfg := config.MustGetConfig(*configFileName)

	// logger config
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.Server.Debug {
		log.Logger = log.With().Caller().Logger()
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		log.Debug().Msg("debug mode is enabled")
	}

	log.Info().Str("config-file", *configFileName).Msg("start mcman")

	// launch minecraft server
	worldReader := world.NewReader(&cfg.Minecraft)

	mcsvr, err := minecraft.NewMinecraftServer(&cfg.Minecraft)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create minecraft server")
	}
	stdout, stderr, err := mcsvr.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start minecraft server")
	}
	log.Info().Msg("started minecraft server")

	// launch log stream
	logStream := logstream.New(&cfg.Minecraft.LogWebhook, map[string]chan string{"stdout": stdout, "stderr": stderr})
	logStream.RegisterLogCallback("webhook", callback.NewWebhookCallback(&cfg.Minecraft.LogWebhook))
	logStream.RegisterLogCallback("zerolog", callback.NewLogCallback(log.With().Str("from", "mc-server").Logger()))
	logStream.Start()

	// launch server
	svr, err := server.New(&cfg.Server, mcsvr, logStream, worldReader)
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

	log.Info().Msg("stopping minecraft server")
	if err := mcsvr.Stop(); err != nil {
		log.Error().Err(err).Msg("cannot stop minecraft server")
	}

	log.Info().Msg("stopping log stream for minecraft server")
	logStream.Stop()
}
