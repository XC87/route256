package product

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"route256.ozon.ru/project/cart/internal/config"
)

const (
	getProductUrl = "/get_product"
	getListSkuUrl = "/list_skus"
)

type ProductService struct {
	client *http.Client
	host   string
	token  string
}

func (service *ProductService) WithTransport(transport Transport) {
	service.client.Transport = &transport
}

func NewProductService(config *config.Config) *ProductService {
	client := &http.Client{
		Timeout: config.ProductServiceTimeout,
	}
	return &ProductService{client: client, host: config.ProductServiceUrl, token: config.ProductServiceToken}
}

var ErrProductNotFound = errors.New("product not found")

func (service *ProductService) GetProduct(sku int64) (*ProductGetProductResponse, error) {
	request := &ProductGetProductRequest{
		service.token,
		sku,
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	resp, err := service.client.Post(service.host+getProductUrl, "application/json", bytes.NewReader(jsonRequest))
	if err != nil {
		return nil, fmt.Errorf("cant connect to product service: %w", err)
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

	return &response, nil
}
