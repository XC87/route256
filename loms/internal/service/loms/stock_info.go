package order_usecase

import (
	"context"
)

func (s *Service) StockInfo(ctx context.Context, sku uint32) (uint64, error) {
	return s.StockRepository.GetCountBySku(ctx, sku)
}
