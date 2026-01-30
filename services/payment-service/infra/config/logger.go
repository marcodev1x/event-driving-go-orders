package config

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Logger() *zap.SugaredLogger {
	logger, err := zap.NewProduction()

	if err != nil {
		panic(err)
	}

	config := zap.NewProductionConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	defer logger.Sync()

	return logger.Sugar()
}
