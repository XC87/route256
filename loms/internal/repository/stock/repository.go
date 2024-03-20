package stock_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"route256.ozon.ru/project/loms/internal/model"
	"sync"
)

type StockMemoryRepository struct {
	sync.RWMutex
	StockData
}

type StockData struct {
	Items map[uint32]Product `json:"items"`
}
type JsonProducts []Product

type Product struct {
	SKU        uint32 `json:"sku"`
	TotalCount uint64 `json:"total_count"`
	Reserved   uint64 `json:"reserved"`
}

func NewStockMemoryRepository() *StockMemoryRepository {
	repository := &StockMemoryRepository{}
	err := repository.GetInitialStockData()
	if err != nil {
		panic(err)
	}
	return repository
}

func (inv *StockMemoryRepository) GetInitialStockData() error {
	data := JsonProducts{}

	file, err := os.ReadFile("stock-data.json")
	if err != nil {
		return fmt.Errorf("error reading stock-data.json: %v", err)
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return fmt.Errorf("error unmarshaling stock data: %v", err)
	}
	inv.Items = make(map[uint32]Product)
	for _, val := range data {
		inv.Items[val.SKU] = val
	}

	return nil
}

/*
	func (inv *StockMemoryRepository) AddItem(ctx context.Context, sku uint32, count uint64) {
		inv.Lock()
		defer inv.Unlock()

		if item, ok := inv.Items[sku]; ok {
			item.TotalCount += count
			inv.Items[sku] = item
		} else {
			inv.Items[sku] = Product{TotalCount: count}
		}
	}
*/
func (inv *StockMemoryRepository) Reserve(ctx context.Context, items []model.Item) error {
	inv.Lock()
	defer inv.Unlock()
	for _, prod := range items {
		item := inv.Items[prod.SKU]
		item.Reserved += prod.Count
		inv.Items[item.SKU] = item
	}

	return nil
}
func (inv *StockMemoryRepository) UnReserve(ctx context.Context, items []model.Item) error {
	inv.Lock()
	defer inv.Unlock()

	for _, prod := range items {
		item := inv.Items[prod.SKU]
		item.Reserved -= prod.Count
		inv.Items[item.SKU] = item
	}

	return nil
}

func (inv *StockMemoryRepository) GetCountBySku(ctx context.Context, sku uint32) (uint64, error) {
	inv.RLock()
	defer inv.RUnlock()
	return inv.Items[sku].TotalCount - inv.Items[sku].Reserved, nil
}

func (inv *StockMemoryRepository) GetBySku(ctx context.Context, sku uint32) (Product, error) {
	inv.RLock()
	defer inv.RUnlock()
	return inv.Items[sku], nil
}
