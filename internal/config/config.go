package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	LogLevel string         `mapstructure:"log_level"`
	RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
	MongoDB  MongoDBConfig  `mapstructure:"mongodb"`
}

type RabbitMQConfig struct {
	URI           string `mapstructure:"uri"`
	Queue         string `mapstructure:"queue"`
	Exchange      string `mapstructure:"exchange"`
	RoutingKey    string `mapstructure:"routing_key"`
	PrefetchCount int    `mapstructure:"prefetch_count"`
}

type MongoDBConfig struct {
	URI        string `mapstructure:"uri"`
	Database   string `mapstructure:"database"`
	Collection string `mapstructure:"collection"`
}

func Load() (*Config, error) {
	v := viper.New()

	v.SetDefault("log_level", "info")
	v.SetDefault("rabbitmq.uri", "amqp://guest:guest@localhost:5672/")
	v.SetDefault("rabbitmq.queue", "events")
	v.SetDefault("rabbitmq.exchange", "events")
	v.SetDefault("rabbitmq.routing_key", "event.#")
	v.SetDefault("rabbitmq.prefetch_count", 10)
	v.SetDefault("mongodb.uri", "mongodb://localhost:27017")
	v.SetDefault("mongodb.database", "cdaq")
	v.SetDefault("mongodb.collection", "events")

	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		var notFound viper.ConfigFileNotFoundError
		if !errors.As(err, &notFound) {
			return nil, fmt.Errorf("reading config file: %w", err)
		}
	}

	v.SetEnvPrefix("CDAQ")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshaling config: %w", err)
	}

	return &cfg, nil
}
