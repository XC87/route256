package service

import (
	"context"
	"sort"
)

func (cartService *CartService) GetItemsByUserId(ctx context.Context, userId int64) (*CartResponse, error) {
	if userId <= 0 {
		return nil, ErrUserInvalid
	}

	skuMap, err := cartService.repository.GetItemsByUserId(userId)
	if err != nil {
		return nil, err
	}

	var cartResponse CartResponse
	if len(skuMap) == 0 {
		return &cartResponse, nil
	}
	skuIdList := make([]int64, 0, len(skuMap))
	for skuId := range skuMap {
		skuIdList = append(skuIdList, skuId)
	}
	sort.Slice(skuIdList, func(i, j int) bool {
		return skuIdList[i] < skuIdList[j]
	})

	list, err := cartService.productService.GetProductList(ctx, skuIdList)
	if err != nil {
		return &cartResponse, err
	}
	for _, productResponse := range list {
		skuId := productResponse.Sku
		item := &CartItem{
			SkuId: skuId,
			Name:  productResponse.Name,
			Count: skuMap[skuId].Count,
			Price: productResponse.Price,
		}

		cartResponse.Items = append(cartResponse.Items, *item)
		cartResponse.TotalPrice += uint32(skuMap[skuId].Count) * productResponse.Price
	}

	return &cartResponse, nil
}
