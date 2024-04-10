package product

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"route256.ozon.ru/project/cart/internal/config"
	"route256.ozon.ru/project/cart/internal/domain"
	"route256.ozon.ru/project/cart/pkg/errgroup"
)

const (
	getProductUrl = "/get_product"
	getListSkuUrl = "/list_skus"
)

type ProductService struct {
	client *http.Client
	host   string
	token  string
	limit  int
}

func (service *ProductService) WithTransport(transport Transport) {
	service.client.Transport = &transport
}

func (service *ProductService) WithRoundTripper(roundTripper http.RoundTripper) {
	service.client.Transport = roundTripper
}

func NewProductService(config *config.Config) *ProductService {
	client := &http.Client{
		Timeout: config.ProductServiceTimeout,
	}
	return &ProductService{client: client,
		host:  config.ProductServiceUrl,
		token: config.ProductServiceToken,
		limit: config.ProductServiceLimit}
}

var ErrProductNotFound = errors.New("product not found")

func (service *ProductService) GetProductList(ctx context.Context, skuList []int64) ([]*domain.Product, error) {
	var productList []*domain.Product

	eg, ctx := errgroup.WithContext(context.Background())
	eg.SetLimit(service.limit)
	// запускаем горуты с отменой запуска если в процессе что-то пошло не так
	for _, sku := range skuList {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("GetProductList innteruped by: %w", ctx.Err())
		default:
			eg.Go(func() (result any, err error) {
				return service.GetProduct(ctx, sku)
			})
		}
	}

	outPutChannel := eg.GetOutChan() // получаем канал
	if outPutChannel != nil {
		for productInfo := range outPutChannel {
			product := productInfo.(*domain.Product)
			productList = append(productList, product)
		}
	}

	if err := eg.Wait(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return productList, nil
}

func (service *ProductService) GetProduct(ctx context.Context, sku int64) (*domain.Product, error) {
	request := &ProductGetProductRequest{
		service.token,
		sku,
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", service.host+getProductUrl, bytes.NewReader(jsonRequest))
	if err != nil {
		return nil, fmt.Errorf("cant connect to product service: %w", err)
	}

	resp, err := service.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cant done request: %w", err)
	}

	defer func() {
		resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrProductNotFound
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response ProductGetProductResponse
	if err = json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	productInfo := domain.Product{
		Name:  response.Name,
		Sku:   sku,
		Price: response.Price,
	}

	return &productInfo, nil
}
