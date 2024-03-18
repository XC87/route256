package service

import "route256.ozon.ru/project/cart/internal/domain"

func (cartService *CartService) AddItem(userId int64, item domain.Item) error {
	if userId <= 0 {
		return ErrUserInvalid
	}
	if item.Count <= 0 {
		return ErrProductCountInvalid
	}

	_, err := cartService.productService.GetProduct(item.Sku_id)
	if err != nil {
		return ErrProductNotFound
	}

	err = cartService.repository.AddItem(userId, item)
	if err != nil {
		return err
	}
	return nil
}
