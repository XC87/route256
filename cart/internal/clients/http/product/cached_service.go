package product

import (
	"context"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"route256.ozon.ru/pkg/keymutex"
	"route256.ozon.ru/project/cart/internal/domain"
	"strconv"
	"strings"
)

type CachedService struct {
	cacher         Cacher
	productService ProductService2
	keyMutex       *keymutex.KeyRWMutex
}

type Cacher interface {
	Get(ctx context.Context, key string) (value string, err error)
	Set(key string, value string) error
}

type ProductService2 interface {
	GetProductList(ctx context.Context, skus []int64) ([]*domain.Product, error)
	GetProduct(ctx context.Context, sku int64) (*domain.Product, error)
}

func NewCachedService(client Cacher, service ProductService2) *CachedService {

	return &CachedService{cacher: client, productService: service, keyMutex: keymutex.NewKeyRWMutex()}
}

func (service *CachedService) GetProduct(ctx context.Context, sku int64) (result *domain.Product, err error) {
	key := strconv.FormatInt(sku, 10)
	service.keyMutex.Lock(key)
	defer service.keyMutex.Unlock(key)
	if product, err := service.cacher.Get(ctx, key); err == nil && product != "" {
		err = json.Unmarshal([]byte(product), &result)
		return result, err
	}

	result, err = service.productService.GetProduct(ctx, sku)
	if err == nil {
		bytes, err := json.Marshal(result)
		err = service.cacher.Set(key, string(bytes))
		if err != nil {
			zap.L().Error("cant set cache in GetProduct", zap.Error(err))
		}
	}

	return result, err
}

func (service *CachedService) GetProductList(ctx context.Context, skus []int64) ([]*domain.Product, error) {
	key := strings.Join(strings.Fields(fmt.Sprint(skus)), ",")
	if productList, err := service.cacher.Get(ctx, key); err == nil && productList != "" {
		var cachedResult []*domain.Product
		err := json.Unmarshal([]byte(productList), &cachedResult)
		return cachedResult, err
	}

	result, err := service.productService.GetProductList(ctx, skus)
	if err == nil {
		bytes, err := json.Marshal(result)
		err = service.cacher.Set(key, string(bytes))
		if err != nil {
			zap.L().Error("cant set cache in GetProduct", zap.Error(err))
		}
	}

	return result, err
}
