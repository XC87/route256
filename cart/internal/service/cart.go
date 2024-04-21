package service

import (
	"context"
	"errors"
	product "route256.ozon.ru/project/cart/internal/clients/http/product"
	"route256.ozon.ru/project/cart/internal/domain"
)

type CartService struct {
	repository     Repository
	productService ProductService
	lomsService    LomsService
}

type LomsService interface {
	GetStockInfo(ctx context.Context, sku uint32) (uint64, error)
	CreateOrder(ctx context.Context, userId int64, items []domain.Item) (int64, error)
}

type Repository interface {
	AddItem(ctx context.Context, userId int64, item domain.Item) error
	DeleteItem(ctx context.Context, userId int64, skuId int64) error
	DeleteItemsByUserId(ctx context.Context, userId int64) error
	GetItemsByUserId(ctx context.Context, userId int64) (domain.ItemsMap, error)
}

type ProductService interface {
	GetProduct(ctx context.Context, sku int64) (*domain.Product, error)
	GetProductList(ctx context.Context, sku []int64) ([]*domain.Product, error)
}

var (
	ErrProductNotFound     = product.ErrProductNotFound
	ErrProductCountInvalid = errors.New("item count invalid")
	ErrStockInsufficient   = errors.New("insufficient stock")
	ErrUserInvalid         = errors.New("user invalid")
	ErrCartIsEmpty         = errors.New("cart is empty")
	ErrCartCantClear       = errors.New("cant clear cart")
	ErrCartCantGet         = errors.New("error fetching users cart")
	ErrOrderCreate         = errors.New("cant create order")
)

func NewCartService(
	repository Repository,
	productService ProductService,
	lomsService LomsService,
) *CartService {
	return &CartService{
		repository:     repository,
		productService: productService,
		lomsService:    lomsService,
	}
}
