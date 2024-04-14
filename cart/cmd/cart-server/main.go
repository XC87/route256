package main

import (
	"context"
	"go.uber.org/zap"
	"log"
	"net/http"
	//_ "net/http/pprof"
	"os"
	"os/signal"
	"route256.ozon.ru/pkg/logger"
	"route256.ozon.ru/pkg/metrics"
	"route256.ozon.ru/pkg/tracer"
	"route256.ozon.ru/project/cart/internal/clients/grpc/loms"
	product "route256.ozon.ru/project/cart/internal/clients/http/product"
	"route256.ozon.ru/project/cart/internal/config"
	"route256.ozon.ru/project/cart/internal/handlers"
	"route256.ozon.ru/project/cart/internal/repository"
	"route256.ozon.ru/project/cart/internal/server"
	"route256.ozon.ru/project/cart/internal/service"
	"syscall"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	cartConfig, err := config.GetConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}

	loggerShutdown := logger.InitLogger(cartConfig.LogLevel, "cart")
	defer loggerShutdown()

	metrics.InitMetricsServer(nil)
	tracerShutdown, err := tracer.InitTracer(cartConfig.TracerUrl, "cart")
	if err == nil {
		defer tracerShutdown(ctx)
	}

	productService := product.NewProductService(cartConfig)
	productService.WithTransport(product.Transport{
		Transport:  http.DefaultTransport,
		RetryCodes: cartConfig.ProductServiceRetryStatus,
		MaxRetries: cartConfig.ProductServiceRetryCount,
	})

	memoryRepository := repository.NewProxyRepository(repository.NewMemoryRepository())
	lomsService, err := loms.NewLomsGrpcClient(ctx, cartConfig.LomsGrpcHost)
	if err != nil {
		zap.L().Fatal("loms grpc client error: ", zap.Error(err))
		return
	}

	cartService := service.NewCartService(memoryRepository, productService, lomsService)

	cartHandler := handlers.NewCartHandler(cartService)
	cartHandler.Register("cart")
	httpServer := server.NewHTTPServer(cartConfig)
	httpServer.Listen(ctx)
}
