package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	LomsGrpcPort string `env:"LOMS_GRPC_PORT, default=:50051"`
	LomsHttpPort string `env:"LOMS_HTTP_PORT, default=:8081"`
}

func GetConfig(ctx context.Context) (*Config, error) {
	var config Config
	if err := envconfig.Process(ctx, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
