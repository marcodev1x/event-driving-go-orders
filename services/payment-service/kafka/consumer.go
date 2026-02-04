package kafka

import (
	"context"
	"errors"
	"payment-service/infra/config"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	handler EventHandler

	numWorkers int
	wg         sync.WaitGroup
}

type EventHandler interface {
	Handle(ctx context.Context, message []byte) error
}

func NewConsumer(brokers []string, topic string, groupID string, handler EventHandler, workers int) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        brokers,
			Topic:          topic,
			GroupID:        groupID,
			MinBytes:       10e3,
			MaxBytes:       10e6,
			CommitInterval: time.Second,
			MaxAttempts:    3, // Número máximo de retrys
			StartOffset:    kafka.LastOffset,
			ErrorLogger:    kafka.LoggerFunc(config.Logger().Errorw),
			MaxWait:        10 * time.Second,
		}),
		handler:    handler,
		numWorkers: workers,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	messageChan := make(chan kafka.Message, c.numWorkers*2)

	for i := 0; i < c.numWorkers; i++ {
		c.wg.Add(1)
		go c.worker(ctx, i, messageChan, &c.wg)
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
	c.wg.Wait()
	return c.reader.Close()
}

func (c *Consumer) worker(ctx context.Context, id int, messages <-chan kafka.Message, wg *sync.WaitGroup) {
	defer wg.Done()

	config.Logger().Infow("Worker iniciado", "worker_id", id)

	for messages := range messages {
		if err := c.handler.Handle(ctx, messages.Value); err != nil {
			config.Logger().Errorw("Erro ao processar mensagem",
				"worker_id", id,
				"error", err,
			)
			continue
		}

		if err := c.reader.CommitMessages(ctx, messages); err != nil {
			config.Logger().Errorw("Erro ao commitar mensagem",
				"worker_id", id,
				"error", err,
			)
		}
	}
}
