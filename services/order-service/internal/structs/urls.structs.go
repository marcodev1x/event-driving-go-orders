package structs

import (
	"order-service/internal/domain"

	"github.com/go-playground/validator/v10"
)

type ById struct {
	Id  string `json:"id" binding:"required"`
	Ref string `json:"ref" binding:"required"`
}

type CreateCheckout struct {
	Price         float64              `json:"price" validate:"required,min=1"`
	PaymentMethod domain.PaymentMethod `json:"payment_method" validate:"required,oneof=boleto credit_card debit_card pix common_money"`
}

func (c *CreateCheckout) ValidateStruct() error {
	validate := validator.New()
	return validate.Struct(c)
}
