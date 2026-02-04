package kafka

import (
	"context"
	"encoding/json"
	"order-service/infra/config"
	"order-service/internal"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"
)

type EventProducer interface {
	PublishEvent(ctx context.Context, key string, event interface{}) error
	PublishEventAsync(ctx context.Context, key string, event interface{}) <-chan error
	Close() error
}

type Producer struct {
	writer *kafka.Writer

	// Para async publishing
	wg sync.WaitGroup
}

func NewProducer(brokers []string, topic string) EventProducer {
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

func (p *Producer) PublishEventAsync(ctx context.Context, key string, event interface{}) <-chan error {
	errCh := make(chan error, 1)

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		defer close(errCh)

		err := p.PublishEvent(ctx, key, event)
		errCh <- err
	}()

	return errCh
}

func (p *Producer) Close() error {
	config.Logger().Info("Aguardando publicações pendentes...")
	p.wg.Wait() // Aguarda todas as goroutines finalizarem

	config.Logger().Info("Fechando producer...")
	return p.writer.Close()
}
