package rest

import (
	"context"
	"encoding/json"
	"notification-service/infra"
	"notification-service/infra/config"
	"notification-service/internal/usecases/interfaces"
	events "notification-service/kafka/events/domain"
	"time"
)

type NotifyConsumer struct {
	usecase interfaces.NotifyImplementation
}

func NewNotifyConsumer(usecases interfaces.NotifyImplementation) *NotifyConsumer {
	return &NotifyConsumer{
		usecase: usecases,
	}
}

func (n *NotifyConsumer) Handle(ctx context.Context, message []byte) error {
	var event events.NotificationInvoice

	if err := json.Unmarshal(message, &event); err != nil {
		config.Logger().Error("Erro ao deserializar evento", err)
		return err
	}

	config.Logger().Info(event)

	switch event.EventType {
	case "payment.confirmed":
		return n.handlePayment(ctx, event)
	default:
		config.Logger().Warnw("Tipo de evento desconhecido", "event_type", event.EventType)
		return nil
	}
}

func (n *NotifyConsumer) handlePayment(ctx context.Context, event events.NotificationInvoice) error {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	defer func() {
		if err := recover(); err != nil {
			config.Logger().Errorw("Erro ao processar evento", "event_id", event.EventID, "error", err)
		}
	}()

	envs := infra.Envs

	config.Logger().Infow("Processando evento de pedido criado",
		"event_id", event.EventID,
		"checkout_id", event.OrderID,
		"stts", event.Status,
		"status", event.Status,
	)

	err := n.usecase.SendNotificationEmail(event, event.ContentID, envs, event.BuyerEmail)

	config.Logger().Infow("Email enviado",
		"event_id", event.EventID,
		"checkout_id", event.OrderID,
		"err", err,
	)

	if err != nil {
		config.Logger().Errorw("Erro ao validar evento", "event_id", event.EventID)
	}

	return nil
}
