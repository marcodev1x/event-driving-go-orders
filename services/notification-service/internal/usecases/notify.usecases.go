package usecases

import (
	"fmt"
	"notification-service/infra/config"
	"notification-service/internal"
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

func (u *NotifyUsecase) SendNotificationEmail(params events.NotificationInvoice, orderId int, envs *config.Env) error {

	fmt.Println(params, orderId)

	Mailer := internal.NewMailer(internal.MailerContent{
		To: struct {
			Email string
		}{
			Email: "vowehe9556@azeriom.com",
		},
		Subject: "Assunto do email",
		Body:    "Corpo do email em HTML",
	})

	msg := Mailer.NewMessage()
	transport := Mailer.Dialer(envs)
	Mailer.MailPrepared(msg, envs)

	if err := transport.DialAndSend(msg); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
