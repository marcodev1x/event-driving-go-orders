package infra

import (
	"payment-service/infra/config"
	"payment-service/internal"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Bootstrap struct{}

func (b *Bootstrap) LoadEnv() *config.Env {
	env := config.LoadEnv()

	config.Logger().Info("Environment variables loaded")

	return env
}

func (b *Bootstrap) RunServer() *gin.Engine {
	router := gin.Default()

	config.Logger().Info("Starting server...")

	return router
}

func (b *Bootstrap) SetupDatabase(env *config.Env) *gorm.DB {
	database := &Database{}
	instance, err := database.Connect(env.DatabaseConfig)

	if err != nil {
		config.Logger().Error("Failed to connect to database", zap.Error(err))
		panic("Failed to connect to database")

		return nil
	}

	return instance
}

func (b *Bootstrap) SetupRedis(env *config.Env) *redis.Client {
	instance := internal.Redis(env)

	if instance == nil {
		return nil
	}

	config.Logger().Info("Redis connected")

	return instance
}
