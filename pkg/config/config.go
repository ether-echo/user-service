package config

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	LogLevel     string   `env:"LOG_LEVEL" envDefault:"DEBUG"`
	KafkaBrokers []string `env:"KAFKA_BROKERS" envSeparator:"," envDefault:"kafka:9092"`
	KafkaTopics  []string `env:"KAFKA_TOPICS" envSeparator:"," envDefault:"/start,/support,message"`
	KafkaGroup   string   `env:"KAFKA_GROUP" envDefault:"metrics-processor-group"`
	DebugPort    string   `env:"DEBUG_PORT" envDefault:"8080"`
	DBPassword   string   `env:"DB_PASSWORD" envDefault:"123"`
	DBHost       string   `env:"DB_HOST" envDefault:"db"`
	DBPort       string   `env:"DB_PORT" envDefault:"5432"`
	DBUser       string   `env:"DB_USER" envDefault:"postgres"`
	DBName       string   `env:"DB_NAME" envDefault:"postgres"`
}

func ReadConfig() (*Config, error) {
	config := Config{}

	err := env.Parse(&config)
	if err != nil {
		return nil, fmt.Errorf("read config error: %w", err)
	}

	return &config, err
}
