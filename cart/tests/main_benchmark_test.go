package tests

import (
	"context"
	"io"
	"net/http"
	"route256.ozon.ru/project/cart/internal/clients/http/product"
	"route256.ozon.ru/project/cart/internal/config"
	"strconv"
	"strings"
	"testing"
)

func BenchmarkInProductService(b *testing.B) {
	ctx := context.Background()
	cartConfig, _ := config.GetConfig(ctx)

	cartConfig.ProductServiceLimit = 1
	for i := 1; i <= 5; i++ {
		b.Run("GetProductList: "+strconv.Itoa(cartConfig.ProductServiceLimit)+" threads", func(b *testing.B) {
			productService := product.NewProductService(cartConfig)
			productService.WithRoundTripper(&mockRoundTrip{})

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				productService.GetProductList(ctx, []int64{1076963, 2956315, 4288068, 4679011, 5097510, 5647362, 6245113, 6966051})
			}
		})
		cartConfig.ProductServiceLimit = i * 3
	}
}

type mockRoundTrip struct{}

func (t *mockRoundTrip) RoundTrip(req *http.Request) (*http.Response, error) {
	response := &http.Response{
		Header:     make(http.Header),
		Request:    req,
		StatusCode: http.StatusOK,
	}
	response.Header.Set("Content-Type", "application/json")

	responseBody :=
		`{
	"name": "тест",
	"price": 123
}`
	//time.Sleep(300 * time.Millisecond)
	response.Body = io.NopCloser(strings.NewReader(responseBody))
	return response, nil
}
