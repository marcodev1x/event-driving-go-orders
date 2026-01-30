package kafka

import (
	"context"
	"encoding/json"
	"payment-service/internal"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(brokers...),
			Topic:        topic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireAll,
			Compression:  kafka.Snappy,
			MaxAttempts:  3,
			BatchSize:    100,
			BatchTimeout: 10 * time.Millisecond,
		},
	}
}

func (p *Producer) PublishEvent(ctx context.Context, key string, event interface{}) error {
	eventData, err := json.Marshal(event)

	if err != nil {
		return internal.NewAPIError("Erro ao serializar evento."+err.Error(), 500, 030)
	}

	msg := kafka.Message{
		Key:   []byte(key),
		Value: eventData,
		Time:  time.Now(),
	}

	if err := p.writer.WriteMessages(ctx, msg); err != nil {
		return internal.NewAPIError("Erro ao publicar evento com key: "+string(msg.Key)+". "+err.Error(), 500, 031)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
