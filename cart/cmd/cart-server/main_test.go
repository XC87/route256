package main

import (
	"bytes"
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
			url:            "http://localhost:8080/user/31337/cart",
			body:           []byte(``),
			expectedStatus: http.StatusNotFound,
			expectedBody:   []byte(`{}`),
		},
		{
			name:           "Add service to cart",
			method:         "POST",
			url:            "http://localhost:8080/user/31337/cart/773297411",
			body:           []byte(`{ "count": 10 }`),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(``),
		},
		{
			name:           "Check cart again, expect 773297411 sku",
			method:         "GET",
			url:            "http://localhost:8080/user/31337/cart",
			body:           []byte(``),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"items":[{"sku_id":773297411,"name":"Кроссовки Nike JORDAN","count":10,"price":2202}],"total_price":22020}`),
		},
		{
			name:           "Add another service to cart",
			method:         "POST",
			url:            "http://localhost:8080/user/31337/cart/2958025",
			body:           []byte(`{ "count": 1 }`),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(``),
		},
		{
			name:           "Check cart again, expect 2958025 sku",
			method:         "GET",
			url:            "http://localhost:8080/user/31337/cart",
			body:           []byte(``),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"items":[{"sku_id":2958025,"name":"Roxy Music. Stranded. Remastered Edition","count":1,"price":1028},{"sku_id":773297411,"name":"Кроссовки Nike JORDAN","count":10,"price":2202}],"total_price":23048}`),
		},
		{
			name:           "Delete 2958025 sku",
			method:         "DELETE",
			url:            "http://localhost:8080/user/31337/cart/2958025",
			body:           []byte(`{}`),
			expectedStatus: http.StatusNoContent,
			expectedBody:   []byte(``),
		},
		{
			name:           "Check cart again, expect 773297411 only sku",
			method:         "GET",
			url:            "http://localhost:8080/user/31337/cart",
			body:           []byte(``),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"items":[{"sku_id":773297411,"name":"Кроссовки Nike JORDAN","count":10,"price":2202}],"total_price":22020}`),
		},
		{
			name:           "Clear cart",
			method:         "DELETE",
			url:            "http://localhost:8080/user/31337/cart",
			body:           []byte(``),
			expectedStatus: http.StatusNoContent,
			expectedBody:   []byte(``),
		}, /*
			{
				name:           "check cart state, expect empty cart",
				method:         "GET",
				url:            "http://localhost:8080/user/31337/cart",
				body:           []byte(``),
				expectedStatus: http.StatusNotFound,
				expectedBody:   []byte(``),
			},
			{
				name:           "Clear cart",
				method:         "DELETE",
				url:            "http://localhost:8080/user/31337/cart",
				body:           []byte(``),
				expectedStatus: http.StatusNoContent,
				expectedBody:   []byte(``),
			},*/
		{
			name:           "Add unknown service, expect error",
			method:         "POST",
			url:            "http://localhost:8080/user/31337/cart/404",
			body:           []byte(`{ "count": 1 }`),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   []byte(``),
		},
	}
	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest(tc.method, tc.url, bytes.NewBuffer(tc.body))
			if err != nil {
				t.Fatalf("error creating request: %v", err)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				t.Fatalf("error sending request: %v", err)
			}
			defer resp.Body.Close()

			body, _ := io.ReadAll(resp.Body)
			if string(body) != string(tc.expectedBody) {
				t.Errorf("unexpected body: want '%s', got '%s'", string(tc.expectedBody), string(body))
			}

			if resp.StatusCode != tc.expectedStatus {
				t.Errorf("unexpected status code: want %d, got %d", tc.expectedStatus, resp.StatusCode)
			}
		})
	}
}
