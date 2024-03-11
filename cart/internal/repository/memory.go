package repository

import (
	"route256.ozon.ru/project/cart/internal/domain"
	"sync"
)

type (
	Memory struct {
		cart map[int64]domain.ItemsMap
		mu   sync.RWMutex
	}
)

func NewMemoryRepository() *Memory {
	return &Memory{
		cart: make(map[int64]domain.ItemsMap),
	}
}

func (memory *Memory) AddItem(userId int64, item domain.Item) error {
	memory.mu.Lock()
	defer memory.mu.Unlock()

	if memory.cart[userId] == nil {
		memory.cart[userId] = domain.ItemsMap{}
	}
	memory.cart[userId][item.Sku_id] = domain.Item{
		Sku_id: item.Sku_id,
		Count:  memory.cart[userId][item.Sku_id].Count + item.Count,
	}

	return nil
}

func (memory *Memory) DeleteItem(userId int64, skuId int64) error {
	memory.mu.Lock()
	defer memory.mu.Unlock()

	delete(memory.cart[userId], skuId)

	return nil
}

func (memory *Memory) DeleteItemsByUserId(userId int64) error {
	memory.mu.Lock()
	defer memory.mu.Unlock()

	delete(memory.cart, userId)

	return nil
}

func (memory *Memory) GetItemsByUserId(userId int64) (domain.ItemsMap, error) {
	memory.mu.RLock()
	defer memory.mu.RUnlock()

	return memory.cart[userId], nil
}
