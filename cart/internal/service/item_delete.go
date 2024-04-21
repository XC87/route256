package service

import (
	"context"
	"go.opentelemetry.io/otel"
)

func (cartService *CartService) DeleteItem(ctx context.Context, userId int64, skuId int64) error {
	_, span := otel.Tracer("default").Start(ctx, "GetItemsByUserId")
	defer span.End()

	if userId <= 0 {
		return ErrUserInvalid
	}

	err := cartService.repository.DeleteItem(ctx, userId, skuId)
	if err != nil {
		return err
	}
	return nil
}
