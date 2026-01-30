package usecases

import (
	"context"
	"payment-service/internal"
	"payment-service/kafka"
	events "payment-service/kafka/events/domain"
	"time"

	"github.com/cenkalti/backoff/v4"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
}

type PaymentUsecase struct {
	// repository mysql.CheckoutImplementation
	cache    Cache
	producer *kafka.Producer
}

func NewPaymentUseCase(cache Cache, producer *kafka.Producer) *PaymentUsecase {
	return &PaymentUsecase{
		cache:    cache,
		producer: producer,
	}
}

func (u *PaymentUsecase) ValidatePayment(params events.OrderCreated, orderId int) error {
	Event := events.PaymentInvoice{
		BaseEvent: events.BaseEvent{
			EventID:   params.EventID,
			ContentID: params.ContentID,
			Timestamp: time.Now(),
		},
		OrderID: orderId,
	}

	bf := backoff.NewExponentialBackOff()
	bf.InitialInterval = 1 * time.Second
	bf.MaxElapsedTime = 10 * time.Second
	bf.MaxInterval = 5 * time.Second
	bf.Multiplier = 2

	if params.Checkout.Price == 100 {
		Event.EventType = "payment.confirmed"

		operationValid := func() error {
			return u.producer.PublishEvent(context.Background(), "payment.confirmed", Event)
		}

		if err := backoff.Retry(operationValid, bf); err != nil {
			return internal.NewAPIError("Erro ao enviar evento mesmo após diversas tentativas. "+err.Error(), 500, 200)
		}
	}

	if params.Checkout.Price == 50 {
		Event.EventType = "payment.failed"

		operationFailed := func() error {
			return u.producer.PublishEvent(context.Background(), "payment.failed", Event)
		}

		if err := backoff.Retry(operationFailed, bf); err != nil {
			return internal.NewAPIError("Erro ao enviar evento mesmo após diversas tentativas. "+err.Error(), 500, 200)
		}
	}

	return nil
}
