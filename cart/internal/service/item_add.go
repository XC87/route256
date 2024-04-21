package service

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"route256.ozon.ru/pkg/logger"
	"route256.ozon.ru/project/cart/internal/domain"
)

func (cartService *CartService) AddItem(ctx context.Context, userId int64, item domain.Item) error {
	ctx, span := otel.Tracer("default").Start(ctx, "AddItem")
	defer span.End()

	if userId <= 0 {
		return ErrUserInvalid
	}
	if item.Count <= 0 {
		return ErrProductCountInvalid
	}

	_, err := cartService.productService.GetProduct(ctx, item.Sku_id)
	if err != nil {
		return ErrProductNotFound
	}

	stockItem, err := cartService.lomsService.GetStockInfo(ctx, uint32(item.Sku_id))
	if err != nil {
		return ErrProductNotFound
	}

	if stockItem < item.Count {
		return ErrStockInsufficient
	}

	err = cartService.repository.AddItem(ctx, userId, item)
	zap.L().With(logger.GetTraceFields(ctx)...).Info("Ctx test log")
	if err != nil {
		return err
	}
	return nil
}
