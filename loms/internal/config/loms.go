package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	LomsGrpcPort              string   `env:"LOMS_GRPC_PORT, default=:50051"`
	LomsHttpPort              string   `env:"LOMS_HTTP_PORT, default=:8081"`
	LomsSharedDbString1       string   `env:"LOMS_DB_SHARDS, default=postgresql://postgres:password@localhost:5432/loms;postgresql://repl_user:repl_password@localhost:5433/loms"`
	LomsSharedDbString2       string   `env:"LOMS_DB_SHARDS, default=postgresql://postgres:password@localhost:5442/loms;postgresql://repl_user:repl_password@localhost:5443/loms"`
	LomsKafkaBrokers          []string `env:"KAFKA_BROKERS, default=localhost:9091,localhost:9092,localhost:9093"`
	LomsKafkaOrderStatusTopic string   `env:"KAFKA_TOPIC_NAME, default=loms.order-events"`
	LogLevel                  string   `env:"LOMS_LOG_LEVEL, default=debug"`
	TracerUrl                 string   `env:"TRACER_URL, default=http://localhost:14268/api/traces"`
}

func GetConfig(ctx context.Context) (*Config, error) {
	var config Config
	if err := envconfig.Process(ctx, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
