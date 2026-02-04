package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	DatabaseConfig *DatabaseConfig
	RedisConfig    *RedisConfig
	KafkaConfig    *KafkaConfig
	SmtpConfig     *SmtpConfig
}

func LoadEnv() *Env {
	err := godotenv.Load()

	getEnv := func(key string) string {
		return os.Getenv(key)
	}

	if err != nil {
		panic(err)
	}

	return &Env{
		DatabaseConfig: &DatabaseConfig{
			Host:     getEnv("DATABASE_HOST"),
			Port:     getEnv("DATABASE_PORT"),
			User:     getEnv("DATABASE_USER"),
			Password: getEnv("DATABASE_PASSWORD"),
			Name:     getEnv("DATABASE_NAME"),
		},
		RedisConfig: &RedisConfig{
			Addr:     getEnv("REDIS_ADDR_URL"),
			Password: getEnv("REDIS_PASSWORD"),
			DB:       getEnv("REDIS_DB"),
			Protocol: getEnv("REDIS_PROTOCOL"),
		},
		KafkaConfig: &KafkaConfig{
			Broker: []string{getEnv("KAFKA_BROKER")},
		},
		SmtpConfig: &SmtpConfig{
			Host:                 getEnv("SMTP_HOST"),
			Port:                 getEnv("SMTP_PORT"),
			Secure:               getEnv("SMTP_SECURE") == "true",
			EmailNoReply:         getEnv("EMAIL_NOREPLY_FROM"),
			EmailNoReplyUser:     getEnv("EMAIL_NOREPLY_USER"),
			EmailNoReplyPassword: getEnv("EMAIL_NOREPLY_PASSWORD"),
		},
	}
}
