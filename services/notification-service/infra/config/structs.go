package config

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       string
	Protocol string
}

type KafkaConfig struct {
	Broker []string
	Topic  string
}
