package order_repository

import (
	"context"
	"errors"
	"sync"

	"route256.ozon.ru/project/loms/internal/model"
)

type OrderMemoryRepository struct {
	orders map[int64]*model.Order
	mu     sync.RWMutex
}

func NewOrderMemoryRepository() *OrderMemoryRepository {
	return &OrderMemoryRepository{
		orders: make(map[int64]*model.Order),
	}
}

var (
	ErrOrderNotFound = errors.New("order not found")
)

func (mr *OrderMemoryRepository) OrderCreate(ctx context.Context, order *model.Order) (int64, error) {
	mr.mu.Lock()
	defer mr.mu.Unlock()

	id := int64(len(mr.orders) + 1)
	order.Id = id
	mr.orders[id] = order

	return id, nil
}

func (mr *OrderMemoryRepository) OrderUpdate(ctx context.Context, order *model.Order) error {
	mr.orders[order.Id] = order
	return nil
}

func (mr *OrderMemoryRepository) OrderInfo(ctx context.Context, id int64) (*model.Order, error) {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	order, ok := mr.orders[id]
	if !ok {
		return nil, ErrOrderNotFound
	}

	return order, nil
}

func (mr *OrderMemoryRepository) OrderPay(ctx context.Context, id int64) error {
	mr.mu.Lock()
	defer mr.mu.Unlock()
	mr.orders[id].Status = model.Paid
	return nil
}

func (mr *OrderMemoryRepository) OrderCancel(ctx context.Context, id int64) error {
	mr.mu.RLock()
	defer mr.mu.RUnlock()

	mr.orders[id].Status = model.Cancelled
	return nil
}
