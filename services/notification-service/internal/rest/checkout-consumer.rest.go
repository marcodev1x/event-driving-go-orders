package rest

import (
	"context"
	"encoding/json"
	"notification-service/infra/config"
	"notification-service/internal/usecases"
	events "notification-service/kafka/events/domain"
)

type CheckoutConsumer struct {
	usecase *usecases.PaymentUsecase
}

func NewCheckoutConsumer(usecases *usecases.PaymentUsecase) *CheckoutConsumer {
	return &CheckoutConsumer{
		usecase: usecases,
	}
}

func (c *CheckoutConsumer) Handle(ctx context.Context, message []byte) error {
	var event events.OrderCreated

	if err := json.Unmarshal(message, &event); err != nil {
		config.Logger().Error("Erro ao deserializar evento", err)
		return err
	}

	switch event.EventType {
	case "order.created":
		return c.handleOrderCreated(ctx, event)
	default:
		config.Logger().Warnw("Tipo de evento desconhecido", "event_type", event.EventType)
		return nil
	}
}

func (c *CheckoutConsumer) handleOrderCreated(ctx context.Context, event events.OrderCreated) error {
	config.Logger().Infow("Processando evento de pedido criado",
		"event_id", event.EventID,
		"checkout_id", event.Checkout.ID,
		"price", event.Checkout.Price,
		"status", event.Checkout.Status,
	)

	err := c.usecase.ValidatePayment(event, event.ContentID)

	if err != nil {
		config.Logger().Errorw("Erro ao validar evento", "event_id", event.EventID)
	}

	return nil
}
