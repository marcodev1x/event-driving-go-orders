package rest

import (
	"context"
	"encoding/json"
	"payment-service/infra/config"
	events "payment-service/kafka/events/domain"
)

// CheckoutConsumer implementa a interface kafka.EventHandler
type CheckoutConsumer struct {
	// adicionar dependências como repository, usecase, etc.
}

func NewCheckoutConsumer() *CheckoutConsumer {
	return &CheckoutConsumer{}
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

	// TODO: Lógica bonus? Adicionar alguma coisa? Evento retorno de order.approved?

	return nil
}
