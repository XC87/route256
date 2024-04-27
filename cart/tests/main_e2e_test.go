package tests

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

type APITestCase struct {
	name           string
	method         string
	url            string
	body           []byte
	expectedStatus int
	expectedBody   []byte
}

func TestAPICases(t *testing.T) {
	testCases := []APITestCase{
		{
			name:           "Check empty cart",
			method:         "GET",
			url:            "http://localhost:8080/user/31336/cart",
			body:           []byte(``),
			expectedStatus: http.StatusNotFound,
			expectedBody:   []byte(`{}`),
		},
		{
			name:           "Add product to cart",
			method:         "POST",
			url:            "http://localhost:8080/user/31336/cart/773297411",
			body:           []byte(`{ "count": 10 }`),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(``),
		},
		{
			name:           "Check cart again, expect 773297411 sku",
			method:         "GET",
			url:            "http://localhost:8080/user/31336/cart",
			body:           []byte(``),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"items":[{"sku_id":773297411,"name":"Кроссовки Nike JORDAN","count":10,"price":2202}],"total_price":22020}`),
		},
		{
			name:           "Add another product to cart",
			method:         "POST",
			url:            "http://localhost:8080/user/31336/cart/2958025",
			body:           []byte(`{ "count": 1 }`),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(``),
		},
		{
			name:           "Check cart again, expect 2958025 sku",
			method:         "GET",
			url:            "http://localhost:8080/user/31336/cart",
			body:           []byte(``),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"items":[{"sku_id":2958025,"name":"Roxy Music. Stranded. Remastered Edition","count":1,"price":1028},{"sku_id":773297411,"name":"Кроссовки Nike JORDAN","count":10,"price":2202}],"total_price":23048}`),
		},
		{
			name:           "Delete 2958025 sku",
			method:         "DELETE",
			url:            "http://localhost:8080/user/31336/cart/2958025",
			body:           []byte(`{}`),
			expectedStatus: http.StatusNoContent,
			expectedBody:   []byte(``),
		},
		{
			name:           "Check cart again, expect 773297411 only sku",
			method:         "GET",
			url:            "http://localhost:8080/user/31336/cart",
			body:           []byte(``),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"items":[{"sku_id":773297411,"name":"Кроссовки Nike JORDAN","count":10,"price":2202}],"total_price":22020}`),
		},
		{
			name:           "Order checkout",
			method:         "POST",
			url:            "http://localhost:8080/cart/checkout",
			body:           []byte(`{ "user_id": 31336 }`),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"orderID":1}`),
		},
		{
			name:           "Add unknown product, expect error",
			method:         "POST",
			url:            "http://localhost:8080/user/31336/cart/404",
			body:           []byte(`{ "count": 1 }`),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   []byte(``),
		},
		{
			name:           "Add another product to cart",
			method:         "POST",
			url:            "http://localhost:8080/user/123/cart/2958025",
			body:           []byte(`{ "count": 1 }`),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(``),
		},
		{
			name:           "Clear cart",
			method:         "DELETE",
			url:            "http://localhost:8080/user/123/cart",
			body:           []byte(``),
			expectedStatus: http.StatusNoContent,
			expectedBody:   []byte(``),
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.url, bytes.NewBuffer(tt.body))
			require.NoError(t, err, "error creating request")
			client := &http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err, "error sending request")
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			assert.Equal(t, string(tt.expectedBody), string(body), "unexpected body")
			assert.Equal(t, tt.expectedStatus, resp.StatusCode, "unexpected status code")
		})
	}
}
