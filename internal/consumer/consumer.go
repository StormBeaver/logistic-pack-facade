package consumer

import (
	"context"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/rs/zerolog/log"
	"github.com/stormbeaver/logistic-kw-pack-facade/internal/config"
	"github.com/stormbeaver/logistic-kw-pack-facade/internal/model"
)

func StartConsuming(ctx context.Context, cfg *config.Kafka, wg *sync.WaitGroup) error {
	config := sarama.NewConfig()

	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(cfg.Brokers, cfg.GroupID, config)
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Strs("topic", cfg.Topics).Msg("subscribe")
		subscribe(ctx, cfg.Topics, consumerGroup, cfg.Tick)
	}()
	return nil
}

func subscribe(
	ctx context.Context,
	topic []string,
	group sarama.ConsumerGroup,
	consumeRate time.Duration,
) {
	ticker := time.NewTicker(consumeRate)
	consumer := model.Consumer{}
	for {
		select {
		case <-ticker.C:
			if err := group.Consume(ctx, topic, &consumer); err != nil {
				log.Error().Err(err).Msg("consume fail")
			}
		case <-ctx.Done():
			return
		}
	}
}
