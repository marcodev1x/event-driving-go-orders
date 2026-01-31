package events

import (
	"notification-service/internal/domain"
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
	Checkout domain.Checkout `json:"checkout"`
}
type PaymentInvoice struct {
	BaseEvent
	OrderID int `json:"order_id"`
}
