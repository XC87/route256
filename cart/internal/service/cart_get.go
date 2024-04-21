package service

import (
	"context"
	"go.opentelemetry.io/otel"
	"sort"
)

func (cartService *CartService) GetItemsByUserId(ctx context.Context, userId int64) (*CartResponse, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "GetItemsByUserId")
	defer span.End()

	if userId <= 0 {
		return nil, ErrUserInvalid
	}

	skuMap, err := cartService.repository.GetItemsByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	span.AddEvent("cartService.repository.GetItemsByUserId")

	var cartResponse CartResponse
	if len(skuMap) == 0 {
		return &cartResponse, nil
	}
	skuIdList := make([]int64, 0, len(skuMap))
	for skuId := range skuMap {
		skuIdList = append(skuIdList, skuId)
	}
	list, err := cartService.productService.GetProductList(ctx, skuIdList)
	span.AddEvent("cartService.productService.GetProductList")
	sort.Slice(list, func(i, j int) bool {
		return list[i].Sku < list[j].Sku
	})
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
