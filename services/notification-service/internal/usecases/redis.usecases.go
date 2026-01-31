package usecases

import (
	"notification-service/internal"
	"time"
)

type RedisUseCase interface {
	Set(key string, value string, ttl time.Duration) error
	Get(key string) (string, error)
}

type RedisImplement struct{}

func NewRedisUsecase() RedisUseCase {
	return &RedisImplement{}
}

func (r *RedisImplement) Set(key string, value string, ttl time.Duration) error {
	return internal.RedisClient.Set(internal.RedisCtx, key, value, ttl).Err()
}

func (r *RedisImplement) Get(key string) (string, error) {
	return internal.RedisClient.Get(internal.RedisCtx, key).Result()
}
