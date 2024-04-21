package service

import (
	"context"
	"go.opentelemetry.io/otel"
)

func (cartService *CartService) DeleteItemsByUserId(ctx context.Context, userId int64) error {
	ctx, span := otel.Tracer("default").Start(ctx, "DeleteItemsByUserId")
	defer span.End()
	if userId <= 0 {
		return ErrUserInvalid
	}

	err := cartService.repository.DeleteItemsByUserId(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}
