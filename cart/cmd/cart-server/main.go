package main

import (
	"log"
	"net/http"
	"route256.ozon.ru/project/cart/internal"
	"route256.ozon.ru/project/cart/internal/client/product"
	"route256.ozon.ru/project/cart/internal/handler"
	"route256.ozon.ru/project/cart/internal/repository"
	"route256.ozon.ru/project/cart/internal/service"
)

func main() {
	config := internal.GetConfig()
	memoryRepository := repository.NewMemoryRepository()
	productService := product.NewProductService(config)

	userFilesService := service.NewCartService(memoryRepository, productService)

	cartHandler := handler.NewCartHandler(userFilesService)
	cartHandler.Register()

	log.Println("Starting server")
	if err := http.ListenAndServe(config.CartHost, nil); err != nil {
		log.Fatal(err)
	}
}
