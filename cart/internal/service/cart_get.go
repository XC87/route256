package service

import (
	"sort"
)

func (cartService *CartService) GetItemsByUserId(userId int64) (*CartResponse, error) {
	if userId <= 0 {
		return nil, ErrUserInvalid
	}

	skuMap, err := cartService.repository.GetItemsByUserId(userId)
	if err != nil {
		return nil, err
	}
	skuIdList := make([]int64, 0, len(skuMap))
	for skuId := range skuMap {
		skuIdList = append(skuIdList, skuId)
	}
	sort.Slice(skuIdList, func(i, j int) bool {
		return skuIdList[i] < skuIdList[j]
	})

	var cartResponse CartResponse
	for _, skuId := range skuIdList {
		productResponse, err := cartService.productService.GetProduct(skuId)
		if err != nil {
			return &cartResponse, err
		}
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
