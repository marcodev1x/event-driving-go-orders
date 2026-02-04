package usecases

import (
	"notification-service/infra/config"
	"notification-service/internal"
	"notification-service/internal/domain"
	"notification-service/internal/templates/payment"
	"notification-service/internal/usecases/interfaces"
	"notification-service/kafka"
	events "notification-service/kafka/events/domain"
	"time"
)

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
}

type NotifyUsecase struct {
	// repository mysql.CheckoutImplementation
	cache    Cache
	producer *kafka.Producer
}

func NewNotifyUseCase(cache Cache, producer *kafka.Producer) interfaces.NotifyImplementation {
	return &NotifyUsecase{
		cache:    cache,
		producer: producer,
	}
}

func (u *NotifyUsecase) SendNotificationEmail(params events.NotificationInvoice, orderId int, envs *config.Env, email string) error {
	subject := "Pagamento confirmado"

	if params.Status == domain.Failed {
		subject = "Pagamento falhou"
	}

	Mailer := internal.NewMailer(internal.MailerContent{
		To: struct {
			Email string
		}{
			Email: email,
		},
		Subject: subject,
	})

	htmlParams := payment.PaymentAlert{
		Name: params.Name,
	}

	msg := Mailer.NewMessage()
	transport := Mailer.Dialer(envs)

	if err := Mailer.SetupHtml(htmlParams, "internal/templates/payment/payment_alert.html"); err != nil {
		config.Logger().Error("Erro ao carregar template", err)
		return err
	}

	Mailer.MailPrepared(msg, envs)

	if err := transport.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}
