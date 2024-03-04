package repository

type Memory struct {
	cart map[int64]map[int64]uint16
}

func NewMemoryRepository() *Memory {
	return &Memory{
		cart: make(map[int64]map[int64]uint16),
	}
}
func (memory *Memory) AddItem(userId int64, skuId int64, count uint16) {
	if memory.cart[userId] == nil {
		memory.cart[userId] = make(map[int64]uint16)
	}
	memory.cart[userId][skuId] += count
}

func (memory *Memory) DeleteItem(userId int64, skuId int64) {
	delete(memory.cart[userId], skuId)
}

func (memory *Memory) DeleteItemsByUserId(userId int64) {
	delete(memory.cart, userId)
}

func (memory *Memory) GetItemsByUserId(userId int64) map[int64]uint16 {
	return memory.cart[userId]
}
