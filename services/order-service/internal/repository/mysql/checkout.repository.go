package mysql

import (
	"order-service/internal/domain"
	"order-service/internal/repository/mysql/interfaces"

	"gorm.io/gorm"
)

type CheckoutRepository struct {
	db *gorm.DB
}

func NewCheckoutRepository(db *gorm.DB) interfaces.CheckoutImplementation {
	return &CheckoutRepository{db}
}

func (r *CheckoutRepository) CreateCheckout(checkout *domain.Checkout) (*domain.Checkout, error) {
	var created *domain.Checkout

	if err := r.db.
		Save(&checkout).
		Find(&created).
		Error; err != nil {
		return nil, err
	}

	return created, nil
}
