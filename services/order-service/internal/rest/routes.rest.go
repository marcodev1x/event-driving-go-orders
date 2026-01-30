package rest

import (
	"order-service/infra"
	"order-service/infra/config"
	"order-service/internal"
	"order-service/internal/middlewares"
	"order-service/internal/repository/mysql"
	"order-service/internal/usecases"
	"order-service/kafka"

	"github.com/gin-gonic/gin"
)

func CheckoutRoutes(env *config.Env) *[]internal.RouteHandler {
	producer := kafka.NewProducer(env.KafkaConfig.Broker, "orders-events")

	consumerHandler := NewCheckoutConsumer()
	_ = kafka.NewConsumer(env.KafkaConfig.Broker, "orders-events", "orders-workers", consumerHandler)

	/*go func() {
		for i := 0; i <= 5; i++ {
			config.Logger().Info("Tentando iniciar Kafka consumer...")

			err := consumer.Start(context.Background())
			if err != nil {
				config.Logger().Error("Erro ao iniciar consumer, retry em 5s", err)
				time.Sleep(5 * time.Second)

				if i == 5 {
					panic("Erro ao iniciar consumer")
				}

				continue
			}

			break
		}
	}()*/

	rest := NewCheckoutRest(
		usecases.NewCheckoutUseCase(mysql.NewCheckoutRepository(infra.DomainDatabase),
			usecases.NewRedisUsecase(),
			producer))

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
