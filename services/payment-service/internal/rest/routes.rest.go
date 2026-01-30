package rest

import (
	"context"
	"payment-service/infra/config"
	"payment-service/internal"
	"payment-service/internal/usecases"
	"payment-service/kafka"
	"time"

	"github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
)

func CheckoutRoutes(env *config.Env) *[]internal.RouteHandler {
	producer := kafka.NewProducer(env.KafkaConfig.Broker, "orders-events")

	bf := backoff.NewExponentialBackOff()
	bf.MaxElapsedTime = 60 * time.Second
	bf.MaxInterval = 5 * time.Second

	go func() {
		config.Logger().Info("Starting Kafka consumer...")

		operation := func() error {
			consumer := kafka.NewConsumer(
				env.KafkaConfig.Broker,
				"orders-events",
				"orders-workers",
				NewCheckoutConsumer(),
			)

			return consumer.Start(context.Background())
		}

		if err := backoff.Retry(operation, bf); err != nil {
			config.Logger().Fatal("Kafka consumer failed after retries", zap.Error(err))
		}
	}()

	usecases.NewPaymentUseCase(usecases.NewRedisUsecase(), producer)

	return &[]internal.RouteHandler{}
}
