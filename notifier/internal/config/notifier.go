package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	KafkaBrokers      []string `env:"KAFKA_BROKERS, default=localhost:9091,localhost:9092,localhost:9093"`
	KafkaTopicNames   []string `env:"KAFKA_TOPIC_NAME, default=loms.order-events"`
	KafkaConsumerName string   `env:"KAFKA_CONSUMER_NAME, default=notifier"`
	TracerUrl         string   `env:"TRACER_URL, default=http://localhost:14268/api/traces"`
}

func GetConfig(ctx context.Context) (*Config, error) {
	var config Config
	if err := envconfig.Process(ctx, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
