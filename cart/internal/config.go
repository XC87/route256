package internal

import (
	"context"
	"github.com/sethvargo/go-envconfig"
	"log"
)

type Config struct {
	CartHost                  string `env:"CART_HOST, default=0.0.0.0:8080"`
	ProductServiceUrl         string `env:"PRODUCT_SERVER_URL, default=http://route256.pavl.uk:8080"`
	ProductServiceToken       string `env:"PRODUCT_SERVER_TOKEN, default=testtoken"`
	ProductServiceRetryStatus []int  `env:"PRODUCT_SERVER_RETRY_STATUS, default=420, 429"`
	ProductServiceRetryCount  int    `env:"PRODUCT_SERVER_RETRY_COUNT, default=3"`
}

func GetConfig() *Config {
	var config Config

	ctx := context.Background()
	if err := envconfig.Process(ctx, &config); err != nil {
		log.Fatal(err)
	}

	return &config
}
