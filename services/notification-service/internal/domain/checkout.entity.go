package domain

import "time"

// Status
const (
	Pending  Status = "pending"
	Paid     Status = "paid"
	Failed   Status = "failed"
	Canceled Status = "canceled"
)

// Payment methods
const (
	Boleto      PaymentMethod = "boleto"
	CreditCard  PaymentMethod = "credit_card"
	DebitCard   PaymentMethod = "debit_card"
	Pix         PaymentMethod = "pix"
	CommonMoney PaymentMethod = "common_money"
)

type Status string
type PaymentMethod string

type Checkout struct {
	ID            int           `json:"ID" gorm:"primaryKey"`
	Price         float64       `json:"price"`
	Status        Status        `json:"status"`
	PaymentMethod PaymentMethod `json:"payment_method"`
	CreatedAt     *time.Time    `json:"createdAt"`
	UpdatedAt     *time.Time    `json:"updatedAt"`
}

func (u *Checkout) TableName() string {
	return "checkout"
}
