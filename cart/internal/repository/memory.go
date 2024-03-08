package repository

import "sync"

type Memory struct {
	cart map[int64]map[int64]uint16
	mu   sync.RWMutex
}

func NewMemoryRepository() *Memory {
	return &Memory{
		cart: make(map[int64]map[int64]uint16),
	}
}

func (memory *Memory) AddItem(userId int64, skuId int64, count uint16) {
	memory.mu.Lock()
	defer memory.mu.Unlock()

	if memory.cart[userId] == nil {
		memory.cart[userId] = make(map[int64]uint16)
	}
	memory.cart[userId][skuId] += count
}

func (memory *Memory) DeleteItem(userId int64, skuId int64) {
	memory.mu.Lock()
	defer memory.mu.Unlock()

	delete(memory.cart[userId], skuId)
}

func (memory *Memory) DeleteItemsByUserId(userId int64) {
	memory.mu.Lock()
	defer memory.mu.Unlock()

	delete(memory.cart, userId)
}

func (memory *Memory) GetItemsByUserId(userId int64) map[int64]uint16 {
	memory.mu.RLock()
	defer memory.mu.RUnlock()

	return memory.cart[userId]
}
