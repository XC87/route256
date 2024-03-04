package product

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"route256.ozon.ru/project/cart/internal"
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

func NewProductService(config *internal.Config) *ProductService {
	// тут конечно вопрос, выносить ли это в main?
	client := &http.Client{
		Transport: &RetryTransport{
			Transport:  http.DefaultTransport,
			RetryCodes: config.ProductServiceRetryStatus,
			MaxRetries: config.ProductServiceRetryCount,
		},
	}
	return &ProductService{client: client, host: config.ProductServiceUrl, token: config.ProductServiceToken}
}

var ErrProductNotFound = errors.New("product not found")

func (service *ProductService) GetProduct(sku int64) (*ProductGetProductResponse, error) {
	request := &ProductGetProductRequest{
		&service.token,
		&sku,
	}
	jsonRequest, err := json.Marshal(request)
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	var response ProductGetProductResponse
	if err = json.Unmarshal(data, &response); err != nil {
		log.Fatal(err)
	}

	return &response, nil
}
