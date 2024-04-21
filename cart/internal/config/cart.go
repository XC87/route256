package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"time"
)

type Config struct {
	CartHost                  string        `env:"CART_HOST, default=0.0.0.0:8080"`
	LomsGrpcHost              string        `env:"LOMS_GRPC_HOST, default=0.0.0.0:50051"`
	ProductServiceUrl         string        `env:"PRODUCT_SERVER_URL, default=http://route256.pavl.uk:8080"`
	ProductServiceToken       string        `env:"PRODUCT_SERVER_TOKEN, default=testtoken"`
	ProductServiceRetryStatus []int         `env:"PRODUCT_SERVER_RETRY_STATUS, default=420, 429"`
	ProductServiceRetryCount  int           `env:"PRODUCT_SERVER_RETRY_COUNT, default=3"`
	ProductServiceTimeout     time.Duration `env:"PRODUCT_SERVER_TIMEOUT, default=5s"`
	ProductServiceLimit       int           `env:"PRODUCT_SERVER_LIMIT, default=10"`
	LogLevel                  string        `env:"CART_LOG_LEVEL, default=debug"`
	TracerUrl                 string        `env:"TRACER_URL, default=http://localhost:14268/api/traces"`
}

func GetConfig(ctx context.Context) (*Config, error) {
	var config Config

	if err := envconfig.Process(ctx, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
