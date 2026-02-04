package rest

import (
	"context"
	"notification-service/infra/config"
	"notification-service/internal"
	"notification-service/internal/usecases"
	"notification-service/kafka"
	"time"

	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
)

func CheckoutRoutes(env *config.Env) *[]internal.RouteHandler {
	producer := kafka.NewProducer(env.KafkaConfig.Broker, "payment-events")

	bf := backoff.NewExponentialBackOff()
	bf.MaxElapsedTime = 60 * time.Second
	bf.MaxInterval = 5 * time.Second

	go func() {
		config.Logger().Info("Starting Kafka consumer...")
		handler := NewNotifyConsumer(usecases.NewNotifyUseCase(usecases.NewRedisUsecase(), producer))

		operation := func() error {
			consumer := kafka.NewConsumer(
				env.KafkaConfig.Broker,
				"payment-events",
				"notification-workers",
				handler,
				10,
			)

			return consumer.Start(context.Background())
		}

		if err := backoff.Retry(operation, bf); err != nil {
			config.Logger().Fatal("Kafka consumer failed after retries", zap.Error(err))
		}
	}()

	return &[]internal.RouteHandler{}
}
