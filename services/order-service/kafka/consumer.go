package kafka

import (
	"context"
	"order-service/infra/config"
	"time"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	reader  *kafka.Reader
	handler EventHandler
}

type EventHandler interface {
	Handle(ctx context.Context, message []byte) error
}

func NewConsumer(brokers []string, topic string, groupID string, handler EventHandler) *Consumer {
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
		handler: handler,
	}
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return c.reader.Close()

		default:
			message, err := c.reader.FetchMessage(ctx)

			if err != nil {
				config.Logger().Error("Erro ao ler mensagem do Kafka", err)
				continue
			}

			if err := c.handler.Handle(ctx, message.Value); err != nil {
				config.Logger().Error("Erro ao processar mensagem do Kafka", err)
				continue
			}

			if err := c.reader.CommitMessages(ctx, message); err != nil {
				config.Logger().Error("Erro ao confirmar mensagem do Kafka", err)
				continue
			}
		}
	}
}
