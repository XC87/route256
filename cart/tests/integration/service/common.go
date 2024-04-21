package service

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"log"
	"route256.ozon.ru/project/cart/internal/clients/grpc/loms"
	"route256.ozon.ru/project/cart/internal/clients/http/product"
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
	ctx := context.Background()
	cartConfig, err := config.GetConfig(ctx)
	require.NoError(s.T(), err)
	s.productService = product.NewProductService(cartConfig)
	s.storage = repository.NewMemoryRepository()
	memoryRepository := repository.NewMemoryRepository()
	productService := product.NewProductService(cartConfig)
	lomsService, err := loms.NewLomsGrpcClient(ctx, cartConfig.LomsGrpcHost)
	if err != nil {
		log.Fatal("loms grpc client error: ", err)
		return
	}

	s.cartService = service.NewCartService(memoryRepository, productService, lomsService)
}

func (s *Suit) TestAdd() {
	userId := int64(1)
	item := domain.Item{
		Sku_id: 4679011,
		Count:  10,
	}
	ctx := context.Background()
	err := s.cartService.AddItem(ctx, userId, item)
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

	ctx := context.Background()
	_ = s.cartService.AddItem(ctx, userId, item)
	res, err := s.cartService.GetItemsByUserId(ctx, userId)
	assert.Equal(s.T(), &cartResponse, res)
	s.Require().NoError(err)
}

func (s *Suit) TestDeleteItem() {
	userId := 1
	item := domain.Item{
		Sku_id: 4679011,
	}
	ctx := context.Background()
	err := s.cartService.DeleteItem(ctx, int64(userId), item.Sku_id)
	s.Require().NoError(err)
}
