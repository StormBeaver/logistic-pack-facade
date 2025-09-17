package consumer

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/stormbeaver/logistic-kw-pack-facade/internal/model"
)

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
					log.Fatal(err)
				}
				if ctx.Err() != nil {
					log.Fatal(ctx.Err())
				}
			case <-ctx.Done():
				wg.Done()
				return
			}
		}
	}()
	return nil
}
