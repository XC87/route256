package service

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"route256.ozon.ru/project/cart/internal/clients/product"
	"route256.ozon.ru/project/cart/internal/config"
	"route256.ozon.ru/project/cart/internal/domain"
	"route256.ozon.ru/project/cart/internal/repository"
	"route256.ozon.ru/project/cart/internal/service"
)

type Suit struct {
	suite.Suite
	productService *product.ProductService
	storage        *repository.Memory
	cartService    *service.CartService
}

func (s *Suit) SetupSuite() {
	var err error
	cartConfig, err := config.GetConfig()
	require.NoError(s.T(), err)
	s.productService = product.NewProductService(cartConfig)
	s.storage = repository.NewMemoryRepository()
	memoryRepository := repository.NewMemoryRepository()
	productService := product.NewProductService(cartConfig)

	s.cartService = service.NewCartService(memoryRepository, productService)
}

func (s *Suit) TestAdd() {
	userId := int64(1)
	item := domain.Item{
		Sku_id: 4679011,
		Count:  10,
	}
	err := s.cartService.AddItem(userId, item)
	s.Require().NoError(err)
}

func (s *Suit) TestList() {
	var cartResponse service.CartResponse
	userId := int64(1)
	item := domain.Item{
		Sku_id: 773297411,
		Count:  10,
	}
	cartItem := service.CartItem{
		SkuId: 773297411,
		Name:  "Кроссовки Nike JORDAN",
		Count: 10,
		Price: 2202,
	}
	cartResponse.Items = append(cartResponse.Items, cartItem)
	cartResponse.TotalPrice += uint32(item.Count) * cartItem.Price

	_ = s.cartService.AddItem(userId, item)
	res, err := s.cartService.GetItemsByUserId(userId)
	assert.Equal(s.T(), &cartResponse, res)
	s.Require().NoError(err)
}

func (s *Suit) TestDeleteItem() {
	userId := 1
	item := domain.Item{
		Sku_id: 4679011,
	}
	err := s.cartService.DeleteItem(int64(userId), item.Sku_id)
	s.Require().NoError(err)
}
