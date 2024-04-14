package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	LomsGrpcPort              string   `env:"LOMS_GRPC_PORT, default=:50051"`
	LomsHttpPort              string   `env:"LOMS_HTTP_PORT, default=:8081"`
	LomsDbUser                string   `env:"POSTGRES_USER, default=postgres"`
	LomsDbPass                string   `env:"POSTGRES_PASSWORD, default=password"`
	LomsDbName                string   `env:"POSTGRES_DB, default=loms"`
	LomsDbSlaveUser           string   `env:"POSTGRESQL_REPLICATION_USER, default=repl_user"`
	LomsDbSlavePass           string   `env:"POSTGRESQL_REPLICATION_PASSWORD, default=repl_password"`
	LomsDbHost                string   `env:"POSTGRES_DB_HOST, default=localhost:5442"`
	LomsDbSlaveHost           string   `env:"POSTGRES_DB_SLAVE_HOST, default=localhost:5433"`
	LomsKafkaBrokers          []string `env:"KAFKA_BROKERS, default=localhost:9091,localhost:9092,localhost:9093"`
	LomsKafkaOrderStatusTopic string   `env:"KAFKA_TOPIC_NAME, default=loms.order-events"`
}

func GetConfig(ctx context.Context) (*Config, error) {
	var config Config
	if err := envconfig.Process(ctx, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
