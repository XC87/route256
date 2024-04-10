package tests

// Всё что ниже для меня, в задании этого не было, мне так проще тестировать

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
)

type httpResponse struct {
	body   []byte
	status int
}

func TestAPICases(t *testing.T) {
	testCases := []struct {
		name           string
		method         string
		url            string
		headers        string
		body           []byte
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name:           "Stock check before create order",
			method:         "POST",
			headers:        "x-auth: test;",
			url:            "http://localhost:8081/stock/info",
			body:           []byte(`{"sku": 773297411}]}`),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"count":"150"}`),
		},
		{
			name:           "Check create order",
			method:         "POST",
			headers:        "x-auth: test;",
			url:            "http://localhost:8081/order/create",
			body:           []byte(`{"user": 123, "items": [{"sku": 773297411, "count": 2}]}`),
			expectedStatus: http.StatusOK,
			//expectedBody:   []byte(`{"orderId":"1"}`),
		},
		{
			name:           "Stock check after create order",
			method:         "POST",
			headers:        "x-auth: test;",
			url:            "http://localhost:8081/stock/info",
			body:           []byte(`{"sku": 773297411}]}`),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"count":"148"}`),
		},
		{
			name:           "Check created order",
			method:         "POST",
			headers:        "x-auth: test;",
			url:            "http://localhost:8081/order/info",
			body:           []byte(`{"orderId": 1}`),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"id":"1", "status":"awaiting_payment", "user":"123", "items":[{"sku":773297411, "count":"2"}]}`),
		},

		{
			name:           "Pay order",
			method:         "POST",
			headers:        "x-auth: test;",
			url:            "http://localhost:8081/order/pay",
			body:           []byte(`{"orderId": 1}`),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{}`),
		},
		{
			name:           "Check payed order",
			method:         "POST",
			headers:        "x-auth: test;",
			url:            "http://localhost:8081/order/info",
			body:           []byte(`{"orderId": 1}`),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"id":"1", "status":"paid", "user":"123", "items":[{"sku":773297411, "count":"2"}]}`),
		},
		{
			name:           "Pay order twice error",
			method:         "POST",
			headers:        "x-auth: test;",
			url:            "http://localhost:8081/order/pay",
			body:           []byte(`{"orderId": 1}`),
			expectedStatus: http.StatusBadRequest,
			//expectedBody:   []byte(`{}`),
		},
		{
			name:           "Check cancel order",
			method:         "POST",
			headers:        "x-auth: test;",
			url:            "http://localhost:8081/order/cancel",
			body:           []byte(`{"orderId": 1}`),
			expectedStatus: http.StatusOK,
			//	expectedBody:   []byte(`{"code":9, "message":"cant reserve", "details":[]}`),
		},
		{
			name:           "Check cancel order twice error",
			method:         "POST",
			headers:        "x-auth: test;",
			url:            "http://localhost:8081/order/cancel",
			body:           []byte(`{"orderId": 1}`),
			expectedStatus: http.StatusBadRequest,
			//		expectedBody:   []byte(`{"code":9, "message":"cant reserve", "details":[]}`),
		},
		{
			name:           "Stock check after cancel order",
			method:         "POST",
			headers:        "x-auth: test;",
			url:            "http://localhost:8081/stock/info",
			body:           []byte(`{"sku": 773297411}]}`),
			expectedStatus: http.StatusOK,
			expectedBody:   []byte(`{"count":"150"}`),
		},
		{
			name:           "Check create order failed",
			method:         "POST",
			headers:        "x-auth: test;",
			url:            "http://localhost:8081/order/create",
			body:           []byte(`{"user": 123, "items": [{"sku": 773292741, "count": 2}]}`),
			expectedStatus: http.StatusBadRequest,
			//	expectedBody:   []byte(`{"code":9, "message":"cant reserve", "details":[]}`),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			resp := sendHTTPRequest(t, tt.method, tt.url, tt.headers, tt.body)
			if tt.expectedBody != nil {
				assert.Equal(t, string(tt.expectedBody), string(resp.body), "unexpected body")
			}
			//	assert.Equal(t, string(tt.expectedBody), string(resp.body), "unexpected body")
			assert.Equal(t, tt.expectedStatus, resp.status, "unexpected status code")
		})
	}
}

func sendHTTPRequest(t *testing.T, method, url string, headers string, body []byte) *httpResponse {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	require.NoError(t, err, "error creating request")

	for _, line := range strings.Split(headers, ";") {
		parts := strings.SplitN(line, ":", 2)
		if len(parts) == 2 {
			key, value := parts[0], parts[1]
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err, "error sending request")
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	return &httpResponse{body: respBody, status: resp.StatusCode}
}
