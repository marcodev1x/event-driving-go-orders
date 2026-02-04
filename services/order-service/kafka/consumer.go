package kafka

import (
	"context"
	"errors"
	"order-service/infra/config"
	"sync"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader     *kafka.Reader
	handler    EventHandler
	numWorkers int
}

type EventHandler interface {
	Handle(ctx context.Context, message []byte) error
}

func NewConsumer(brokers []string, topic string, groupID string, handler EventHandler) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:  brokers,
			Topic:    topic,
			GroupID:  groupID,
			MinBytes: 10e3,
			MaxBytes: 10e6,
		}),
		handler:    handler,
		numWorkers: 5,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	var wg sync.WaitGroup

	messageChan := make(chan kafka.Message, c.numWorkers*2)

	for i := 0; i < c.numWorkers; i++ {
		wg.Add(1)
		go c.worker(ctx, i, messageChan, &wg)
	}

	go func() {
		defer close(messageChan)

		for {
			select {
			case <-ctx.Done():
				return
			default:
				msg, err := c.reader.FetchMessage(ctx)
				if err != nil {
					if errors.Is(err, context.Canceled) {
						return
					}
					config.Logger().Error("Erro ao ler mensagem", err)
					continue
				}

				messageChan <- msg
			}
		}
	}()

	<-ctx.Done()
	wg.Wait()
	return c.reader.Close()
}

// Worker processa mensagens do canal
func (c *Consumer) worker(ctx context.Context, id int, messages <-chan kafka.Message, wg *sync.WaitGroup) {
	defer wg.Done()

	config.Logger().Infow("Worker iniciado", "worker_id", id)

	for msg := range messages {
		if err := c.handler.Handle(ctx, msg.Value); err != nil {
			config.Logger().Errorw("Erro ao processar mensagem",
				"worker_id", id,
				"error", err,
			)
			continue
		}

		if err := c.reader.CommitMessages(ctx, msg); err != nil {
			config.Logger().Errorw("Erro ao commitar mensagem",
				"worker_id", id,
				"error", err,
			)
		}
	}
}
