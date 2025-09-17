package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stormbeaver/logistic-kw-pack-facade/internal/config"
	"github.com/stormbeaver/logistic-kw-pack-facade/internal/consumer"
)

func main() {

	sigs := make(chan os.Signal, 1)
	wg := &sync.WaitGroup{}

	if err := config.ReadConfigYML("config.yml"); err != nil {
		log.Fatal().Err(err).Msg("Failed init configuration")
	}

	cfg := config.GetConfigInstance()

	log.Info().
		Str("version", cfg.Project.Version).
		Str("commitHash", cfg.Project.CommitHash).
		Bool("debug", cfg.Project.Debug).
		Str("environment", cfg.Project.Environment).
		Msgf("Starting service: %s", cfg.Project.Name)

	if cfg.Project.Debug {
		log.Level(zerolog.DebugLevel)
	} else {
		log.Level(zerolog.InfoLevel)
	}

	ctx, cancel := context.WithCancel(context.Background())

	err := consumer.StartConsuming(ctx, &cfg.Kafka, wg)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
	cancel()
	wg.Wait()
}
