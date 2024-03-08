package service

import (
	"errors"
	"route256.ozon.ru/project/cart/internal/client/product"
	"sort"
)

type CartService struct {
	productService ProductService
	repository     Repository
}

type Repository interface {
	AddItem(userId int64, skuId int64, count uint16)
	DeleteItem(userId int64, skuId int64)
	DeleteItemsByUserId(userId int64)
	GetItemsByUserId(userId int64) map[int64]uint16
}

type ProductService interface {
	GetProduct(sku int64) (*product.ProductGetProductResponse, error)
}

var (
	ErrProductNotFound     = product.ErrProductNotFound
	ErrProductCountInvalid = errors.New("item count invalid")
	ErrUserInvalid         = errors.New("user invalid")
)

func NewCartService(
	repository Repository,
	productService ProductService,
) *CartService {
	return &CartService{
		productService: productService,
		repository:     repository,
	}
}
func (cartService *CartService) AddItem(userId int64, skuId int64, count uint16) error {
	if userId <= 0 {
		return ErrUserInvalid
	}
	if count <= 0 {
		return ErrProductCountInvalid
	}

	_, err := cartService.productService.GetProduct(skuId)
	if err != nil {
		return ErrProductNotFound
	}

	cartService.repository.AddItem(userId, skuId, count)
	return nil
}

func (cartService *CartService) DeleteItem(userId int64, skuId int64) error {
	if userId <= 0 {
		return ErrUserInvalid
	}

	cartService.repository.DeleteItem(userId, skuId)
	return nil
}
func (cartService *CartService) DeleteItemsByUserId(userId int64) error {
	if userId <= 0 {
		return ErrUserInvalid
	}

	cartService.repository.DeleteItemsByUserId(userId)
	return nil
}

func (cartService *CartService) GetItemsByUserId(userId int64) (*CartResponse, error) {
	if userId <= 0 {
		return nil, ErrUserInvalid
	}

	skuMap := cartService.repository.GetItemsByUserId(userId)
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
			SkuID: skuId,
			Name:  productResponse.Name,
			Count: skuMap[skuId],
			Price: productResponse.Price,
		}

		cartResponse.Items = append(cartResponse.Items, *item)
		cartResponse.TotalPrice += uint32(skuMap[skuId]) * productResponse.Price
	}

	return &cartResponse, nil
}
