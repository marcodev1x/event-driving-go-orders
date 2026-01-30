package usecases

import (
	"context"
	"order-service/infra/config"
	"order-service/internal"
	"order-service/internal/domain"
	"order-service/internal/repository/mysql"
	"order-service/internal/structs"
	"order-service/kafka"
	"order-service/kafka/events/domain"

	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
}

type CheckoutUsecase struct {
	repository mysql.CheckoutImplementation
	cache      Cache
	producer   *kafka.Producer
}

func NewCheckoutUseCase(repo mysql.CheckoutImplementation, cache Cache, producer *kafka.Producer) *CheckoutUsecase {
	return &CheckoutUsecase{
		repository: repo,
		cache:      cache,
		producer:   producer,
	}
}

func (u *CheckoutUsecase) CreateCheckout(req structs.CreateCheckout) (bool, error) {
	if err := req.ValidateStruct(); err != nil {
		return false, internal.NewAPIError("Estrutura inválida."+err.Error(), 400, 101)
	}

	checkout := &domain.Checkout{
		Price:         req.Price,
		Status:        domain.Pending,
		PaymentMethod: req.PaymentMethod,
	}

	created, err := u.repository.CreateCheckout(checkout)

	if err != nil {
		return false, internal.NewAPIError("Erro ao criar checkout.", 500, 102)
	}

	event := events.OrderCreated{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New().String(),
			EventType: "order.created",
			ContentID: created.ID,
			Timestamp: time.Now(),
		},
		Checkout: *checkout,
	}

	producerOperation := func() error {
		return u.producer.PublishEvent(context.Background(), event.EventType, event)
	}

	go func() {
		// backoff exponencial se falhar.
		if err := producerOperation(); err != nil {
			retry := backoff.NewExponentialBackOff()
			retry.InitialInterval = 1 * time.Second
			retry.MaxElapsedTime = 10 * time.Second
			retry.MaxInterval = 5 * time.Second
			retry.Multiplier = 2

			if err := backoff.Retry(producerOperation, retry); err != nil {
				config.Logger().Error("Erro ao publicar evento após retries.", err)
			}
		}
	}()

	return true, nil
}
