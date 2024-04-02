package service

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/cart/internal/domain"
)

func (cartService *CartService) OrderCheckout(ctx context.Context, userId int64) (int64, error) {
	if userId <= 0 {
		return 0, ErrUserInvalid
	}

	cart, err := cartService.GetItemsByUserId(ctx, userId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", ErrCartCantGet, err)
	}

	if len(cart.Items) == 0 {
		return 0, ErrCartIsEmpty
	}

	items := convertCartToDomainItems(cart)
	orderId, err := cartService.lomsService.CreateOrder(ctx, userId, items)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", ErrOrderCreate, err)
	}

	err = cartService.DeleteItemsByUserId(userId)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", ErrCartCantClear, err)
	}

	return orderId, nil
}

func convertCartToDomainItems(cart *CartResponse) []domain.Item {
	result := make([]domain.Item, len(cart.Items))
	for i, item := range cart.Items {
		result[i] = domain.Item{
			Sku_id: item.SkuId,
			Count:  item.Count,
		}
	}

	return result
}
