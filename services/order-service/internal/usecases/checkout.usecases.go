package usecases

import (
	"context"
	"order-service/infra/config"
	"order-service/internal"
	"order-service/internal/domain"
	kafkainterfaces "order-service/internal/interfaces/kafka"
	"order-service/internal/repository/mysql/interfaces"
	"order-service/internal/structs"
	redisinterfaces "order-service/internal/usecases/interfaces"
	events "order-service/kafka/events/domain"

	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
)

type CheckoutUsecase struct {
	repository interfaces.CheckoutImplementation
	cache      redisinterfaces.RedisImplementation
	producer   kafkainterfaces.EventProducer
}

func NewCheckoutUseCase(repo interfaces.CheckoutImplementation, cache redisinterfaces.RedisImplementation, producer kafkainterfaces.EventProducer) *CheckoutUsecase {
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

	err := u.repository.CreateCheckout(checkout)

	if err != nil {
		return false, internal.NewAPIError("Erro ao criar checkout.", 500, 102)
	}

	event := events.OrderCreated{
		BaseEvent: events.BaseEvent{
			EventID:   uuid.New().String(),
			EventType: "order.created",
			ContentID: checkout.ID,
			Timestamp: time.Now(),
		},
		Checkout: *checkout,
	}

	go u.publishWithRetry(event)

	return true, nil
}

func (u *CheckoutUsecase) publishWithRetry(event events.OrderCreated) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	operation := func() error {
		return u.producer.PublishEvent(ctx, event.EventType, event)
	}

	// Configuração do retry exponencial
	retry := backoff.NewExponentialBackOff()
	retry.InitialInterval = 1 * time.Second
	retry.MaxInterval = 10 * time.Second
	retry.MaxElapsedTime = 30 * time.Second
	retry.Multiplier = 2

	if err := backoff.Retry(operation, retry); err != nil {
		config.Logger().Errorw("Falha ao publicar evento após múltiplas tentativas",
			"error", err,
			"event_id", event.EventID,
			"event_type", event.EventType,
			"checkout_id", event.Checkout.ID,
		)
		return
	}

	config.Logger().Infow("Evento publicado com sucesso",
		"event_id", event.EventID,
		"event_type", event.EventType,
		"checkout_id", event.Checkout.ID,
	)
}
