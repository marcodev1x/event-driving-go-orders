package events

import (
	"payment-service/internal/domain"
	"time"
)

type BaseEvent struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	ContentID int       `json:"content_id"`
}

type OrderCreated struct {
	BaseEvent
	Checkout   domain.Checkout `json:"checkout"`
	BuyerEmail string          `json:"buyer_email"`
	Name       string          `json:"name"`
}
type PaymentInvoice struct {
	BaseEvent
	Name          string               `json:"name"`
	Status        domain.Status        `json:"status"`
	OrderID       int                  `json:"order_id"`
	BuyerEmail    string               `json:"buyer_email"`
	PaymentMethod domain.PaymentMethod `json:"payment_method"`
	Price         float64              `json:"price"`
}
