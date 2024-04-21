package order_usecase

import (
	"context"
	"go.opentelemetry.io/otel"
)

func (s *Service) StockInfo(ctx context.Context, sku uint32) (uint64, error) {
	_, span := otel.Tracer("default").Start(ctx, "StockInfo")
	defer span.End()

	return s.StockRepository.GetCountBySku(ctx, sku)
}
