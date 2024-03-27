package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"log"
)

type Config struct {
	LomsGrpcPort string `env:"LOMS_GRPC_PORT, default=:50051"`
	LomsHttpPort string `env:"LOMS_HTTP_PORT, default=:8081"`
}

func GetConfig(ctx context.Context) *Config {
	var config Config
	if err := envconfig.Process(ctx, &config); err != nil {
		log.Fatalf("GetConfig: %s", err)
	}

	return &config
}
