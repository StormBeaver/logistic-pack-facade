package consumer

import (
	"context"
	"sync"

	"github.com/IBM/sarama"
	"github.com/stormbeaver/logistic-kw-pack-facade/internal/config"
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
	}

	return nil
}
