package interfaces

import (
	"notification-service/infra/config"
	events "notification-service/kafka/events/domain"
	"time"
)

type RedisImplementation interface {
	Set(key string, value string, ttl time.Duration) error
	Get(key string) (string, error)
}

type NotifyImplementation interface {
	SendNotificationEmail(params events.NotificationInvoice, orderId int, envs *config.Env, email string) error
}
