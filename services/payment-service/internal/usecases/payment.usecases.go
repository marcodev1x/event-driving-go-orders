package usecases

import (
	"payment-service/kafka"
	"time"
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

func (p *PaymentUsecase) ValidatePayment() bool {
	return true
}
