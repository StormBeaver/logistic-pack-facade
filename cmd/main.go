package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"sync"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stormbeaver/logistic-kw-pack-facade/internal/config"
	"github.com/stormbeaver/logistic-kw-pack-facade/internal/consumer"
)

func main() {
	go http.ListenAndServe("localhost:10000", nil)
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
		Strs("brokers", cfg.Kafka.Brokers).
		Strs("topics", cfg.Kafka.Topics).
		Msgf("Starting service: %s", cfg.Project.Name)

	if cfg.Project.Debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	err := consumer.StartConsuming(ctx, &cfg.Kafka, wg)

	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	log.Info().Msg("start consuming")

	<-ctx.Done()
	wg.Wait()
}
