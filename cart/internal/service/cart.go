package service

import (
	"errors"
	"route256.ozon.ru/project/cart/internal/clients/product"
	"route256.ozon.ru/project/cart/internal/domain"
)

type CartService struct {
	productService ProductService
	repository     Repository
}

type Repository interface {
	AddItem(userId int64, item domain.Item) error
	DeleteItem(userId int64, skuId int64) error
	DeleteItemsByUserId(userId int64) error
	GetItemsByUserId(userId int64) (domain.ItemsMap, error)
}

type ProductService interface {
	GetProduct(sku int64) (*product.ProductGetProductResponse, error)
}

var (
	ErrProductNotFound     = product.ErrProductNotFound
	ErrProductCountInvalid = errors.New("item count invalid")
	ErrUserInvalid         = errors.New("user invalid")
)

func NewCartService(
	repository Repository,
	productService ProductService,
) *CartService {
	return &CartService{
		productService: productService,
		repository:     repository,
	}
}
