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

	for _, topic := range cfg.Topics {
		subscribe(ctx, topic, consumerGroup, cfg.Tick, wg)
		log.Debug().Str("topic:", topic).Msg("subscribe to topic")
	}
	return nil
}

func subscribe(
	ctx context.Context,
	topic string,
	group sarama.ConsumerGroup,
	consumeRate time.Duration,
	wg *sync.WaitGroup,
) error {
	ticker := time.NewTicker(consumeRate)
	consumer := model.Consumer{}
	wg.Add(1)
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := group.Consume(ctx, []string{topic}, &consumer); err != nil {
					log.Error().Err(err).Msg("failed to consume")
				}
			case <-ctx.Done():
				wg.Done()
				return
			}
		}
	}()
	return nil
}
