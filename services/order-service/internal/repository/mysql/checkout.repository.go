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

func (r *CheckoutRepository) CreateCheckout(checkout *domain.Checkout) error {
	if err := r.db.
		Save(&checkout).
		Find(&checkout).
		Error; err != nil {
		return err
	}

	return nil
}
