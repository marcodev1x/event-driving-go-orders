package usecases

import (
	"context"
	"payment-service/infra/config"
	"payment-service/internal"
	"payment-service/internal/domain"
	"payment-service/kafka"
	events "payment-service/kafka/events/domain"
	"time"

	"github.com/cenkalti/backoff/v4"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
}

type PaymentImplementation interface {
	ValidatePayment(params events.OrderCreated, orderId int) error
}

type PaymentUsecase struct {
	// repository mysql.CheckoutImplementation
	cache    Cache
	producer kafka.EventProducer
}

func NewPaymentUseCase(cache Cache, producer kafka.EventProducer) PaymentImplementation {
	return &PaymentUsecase{
		cache:    cache,
		producer: producer,
	}
}

func (u *PaymentUsecase) ValidatePayment(params events.OrderCreated, orderId int) error {
	OrderReturnEvent := events.PaymentInvoice{
		BaseEvent: events.BaseEvent{
			EventID:   params.EventID,
			ContentID: params.ContentID,
			Timestamp: time.Now(),
		},
		Name:          params.Name,
		OrderID:       orderId,
		BuyerEmail:    params.BuyerEmail,
		PaymentMethod: params.Checkout.PaymentMethod,
		Price:         params.Checkout.Price,
	}

	bf := backoff.NewExponentialBackOff()
	bf.InitialInterval = 1 * time.Second
	bf.MaxElapsedTime = 10 * time.Second
	bf.MaxInterval = 5 * time.Second
	bf.Multiplier = 2

	if params.Checkout.Price == 100 {
		OrderReturnEvent.EventType = "payment.confirmed"
		OrderReturnEvent.Status = domain.Paid

		err := u.publishWithRetry(OrderReturnEvent)

		if err != nil {
			return err
		}
	}

	if params.Checkout.Price == 50 {
		OrderReturnEvent.EventType = "payment.failed"
		OrderReturnEvent.Status = domain.Failed

		err := u.publishWithRetry(OrderReturnEvent)
		if err != nil {
			return err
		}
	}

	return nil
}

func (p *PaymentUsecase) publishWithRetry(event events.PaymentInvoice) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	operation := func() error {
		return p.producer.PublishEvent(ctx, event.EventType, event)
	}

	retry := backoff.NewExponentialBackOff()
	retry.InitialInterval = 1 * time.Second
	retry.MaxInterval = 5 * time.Second
	retry.MaxElapsedTime = 10 * time.Second
	retry.Multiplier = 2

	if err := backoff.Retry(operation, retry); err != nil {
		config.Logger().Errorw("Falha ao publicar evento após múltiplas tentativas",
			"error", err,
			"event_id", event.EventID,
			"event_type", event.EventType,
			"checkout_id", event.OrderID,
		)
		return internal.NewAPIError("Erro ao enviar evento mesmo após diversas tentativas.", 500, 200)
	}

	config.Logger().Infow("Evento publicado com sucesso",
		"event_id", event.EventID,
		"event_type", event.EventType,
		"checkout_id", event.OrderID,
	)

	return nil
}
