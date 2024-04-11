package service

import (
	"context"
	"route256.ozon.ru/project/cart/internal/domain"
)

func (cartService *CartService) AddItem(ctx context.Context, userId int64, item domain.Item) error {
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

	err = cartService.repository.AddItem(userId, item)
	if err != nil {
		return err
	}
	return nil
}
