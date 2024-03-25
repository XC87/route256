package order_usecase

import (
	"context"
	"errors"
	"route256.ozon.ru/project/loms/internal/model"
	order_repository "route256.ozon.ru/project/loms/internal/repository/pgs/order"
)

type OrderRepository interface {
	OrderCreate(ctx context.Context, order *model.Order) (int64, error)
	OrderUpdate(ctx context.Context, order *model.Order) error
	OrderInfo(ctx context.Context, id int64) (*model.Order, error)
	OrderPay(ctx context.Context, id int64) error
	OrderCancel(ctx context.Context, id int64) error
}

type StockRepository interface {
	GetCountBySku(ctx context.Context, sku uint32) (uint64, error)
	GetBySku(ctx context.Context, sku uint32) (*model.ProductStock, error)
	Reserve(ctx context.Context, items []model.Item) error
	UnReserve(ctx context.Context, items []model.Item) error
}

type Service struct {
	OrderRepository OrderRepository
	StockRepository StockRepository
}

var (
	ErrOrderNotFound        = order_repository.ErrOrderNotFound
	ErrOrderInvalid         = errors.New("invalid order")
	ErrOrderCantReserve     = errors.New("cant reserve")
	ErrOrderCantUnReserve   = errors.New("cant unreserve")
	ErrOrderSkuNotFound     = errors.New("sku not found")
	ErrOrderCantBePaid      = errors.New("order cant be paid")
	ErrOrderAlreadyPaid     = errors.New("order already paid")
	ErrOrderAlreadyCanceled = errors.New("order already canceled")
	ErrStockNotEnough       = errors.New("not enough items in stock")
	ErrReserveNotEnough     = errors.New("not enough items in reserve")
)

func NewService(orderRepository OrderRepository, stockRepository StockRepository) *Service {
	return &Service{OrderRepository: orderRepository, StockRepository: stockRepository}
}

func checkIfCanReserve(rep StockRepository, items []model.Item) error {
	for _, item := range items {
		itemStock, err := rep.GetBySku(context.Background(), item.SKU)
		if err != nil {
			return err
		}

		if itemStock.Reserved+item.Count > itemStock.TotalCount {
			return ErrStockNotEnough
		}
	}

	return nil
}

func checkIfCanUnReserve(rep StockRepository, items []model.Item) error {
	for _, item := range items {
		itemStock, err := rep.GetBySku(context.Background(), item.SKU)
		if err != nil {
			return err
		}

		if itemStock.Reserved < item.Count {
			return ErrReserveNotEnough
		}
	}

	return nil
}
