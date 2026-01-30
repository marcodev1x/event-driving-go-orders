package rest

import (
	"context"
	"order-service/infra"
	"order-service/infra/config"
	"order-service/internal"
	"order-service/internal/middlewares"
	"order-service/internal/repository/mysql"
	"order-service/internal/usecases"
	"order-service/kafka"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func CheckoutRoutes(env *config.Env) *[]internal.RouteHandler {
	producer := kafka.NewProducer(env.KafkaConfig.Broker, "orders-events")

	checkoutUseCase := usecases.NewCheckoutUseCase(
		mysql.NewCheckoutRepository(infra.DomainDatabase),
		usecases.NewRedisUsecase(),
		producer,
	)

	bf := backoff.NewExponentialBackOff()
	bf.MaxElapsedTime = 60 * time.Second
	bf.MaxInterval = 5 * time.Second

	go func() {
		config.Logger().Info("Starting Kafka consumer...")

		operation := func() error {
			consumer := kafka.NewConsumer(
				env.KafkaConfig.Broker,
				"payment-events",
				"order-payment-handler",
				NewCheckoutConsumer(),
			)

			return consumer.Start(context.Background())
		}

		if err := backoff.Retry(operation, bf); err != nil {
			config.Logger().Fatal(
				"Kafka consumer failed after retries",
				zap.Error(err),
			)
		}
	}()

	rest := NewCheckoutRest(checkoutUseCase)

	return &[]internal.RouteHandler{
		{
			Path:    "/create-checkout",
			Handler: rest.CreateCheckout,
			Method:  internal.POST,
			Middlewares: []gin.HandlerFunc{
				middlewares.Interceptors.ErrorHandler(),
				middlewares.Interceptors.RateLimiter(1, 1),
			},
		},
	}
}
