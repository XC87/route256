package main

import (
	"log"
	"net/http"
	"route256.ozon.ru/project/cart/internal/clients/product"
	"route256.ozon.ru/project/cart/internal/config"
	"route256.ozon.ru/project/cart/internal/handlers"
	"route256.ozon.ru/project/cart/internal/repository"
	"route256.ozon.ru/project/cart/internal/service"
)

func main() {
	cartConfig, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	productService := product.NewProductService(cartConfig)
	productService.WithTransport(product.Transport{
		Transport:  http.DefaultTransport,
		RetryCodes: cartConfig.ProductServiceRetryStatus,
		MaxRetries: cartConfig.ProductServiceRetryCount,
	})

	memoryRepository := repository.NewMemoryRepository()
	cartService := service.NewCartService(memoryRepository, productService)

	cartHandler := handlers.NewCartHandler(cartService)
	cartHandler.Register()

	log.Println("Starting server")
	if err = http.ListenAndServe(cartConfig.CartHost, nil); err != nil {
		log.Fatal(err)
	}
}
