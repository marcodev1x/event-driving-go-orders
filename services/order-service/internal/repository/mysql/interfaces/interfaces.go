package interfaces

import (
	"order-service/internal/domain"
)

type CheckoutImplementation interface {
	CreateCheckout(checkout *domain.Checkout) (*domain.Checkout, error)
}
