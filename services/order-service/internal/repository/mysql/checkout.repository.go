package mysql

import (
	"fmt"
	"order-service/internal/domain"

	"gorm.io/gorm"
)

type CheckoutImplementation interface {
	CreateCheckout(checkout *domain.Checkout) (*domain.Checkout, error)
}

type CheckoutRepository struct {
	db *gorm.DB
}

func NewCheckoutRepository(db *gorm.DB) CheckoutImplementation {
	return &CheckoutRepository{db}
}

func (r *CheckoutRepository) CreateCheckout(checkout *domain.Checkout) (*domain.Checkout, error) {
	var created *domain.Checkout

	fmt.Println("Creating checkout", checkout)

	if err := r.db.
		Save(&checkout).
		Find(&created).
		Error; err != nil {
		return nil, err
	}

	return created, nil
}
