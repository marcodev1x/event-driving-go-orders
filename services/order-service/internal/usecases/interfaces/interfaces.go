package interfaces

import "time"

type RedisImplementation interface {
	Set(key string, value string, ttl time.Duration) error
	Get(key string) (string, error)
}
