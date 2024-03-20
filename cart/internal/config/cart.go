package config

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"time"
)

type Config struct {
	CartHost                  string        `env:"CART_HOST, default=0.0.0.0:8080"`
	LomsGrpcHost              string        `env:"LOMS_HOST, default=0.0.0.0:50051"`
	ProductServiceUrl         string        `env:"PRODUCT_SERVER_URL, default=http://route256.pavl.uk:8080"`
	ProductServiceToken       string        `env:"PRODUCT_SERVER_TOKEN, default=testtoken"`
	ProductServiceRetryStatus []int         `env:"PRODUCT_SERVER_RETRY_STATUS, default=420, 429"`
	ProductServiceRetryCount  int           `env:"PRODUCT_SERVER_RETRY_COUNT, default=3"`
	ProductServiceTimeout     time.Duration `env:"PRODUCT_SERVER_TIMEOUT, default=5s"`
}

func GetConfig(ctx context.Context) (*Config, error) {
	var config Config
	if err := envconfig.Process(ctx, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
