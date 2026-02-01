package kafkainterfaces

import "context"

type EventProducer interface {
	PublishEvent(ctx context.Context, key string, payload any) error
}
