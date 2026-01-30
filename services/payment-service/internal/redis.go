package internal

import (
	"context"
	"payment-service/infra/config"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var (
	RedisCtx    context.Context
	RedisClient *redis.Client
)

func Redis(envs *config.Env) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     envs.RedisConfig.Addr,
		Password: envs.RedisConfig.Password,
		DB:       0,
		Protocol: 2,
	})

	RedisCtx = context.Background()

	if err := rdb.Ping(RedisCtx).Err(); err != nil {
		config.Logger().Error("Failed to connect to Redis", zap.Error(err))

		return nil
	}

	RedisClient = rdb

	return rdb
}
