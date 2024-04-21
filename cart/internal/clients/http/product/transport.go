package product

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

type Transport struct {
	Transport  http.RoundTripper
	RetryCodes []int
	MaxRetries int
}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error
	var bodyBytes []byte

	// копия чтобы "перематывать"
	if req.Body != nil {
		bodyBytes, err = io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
	}

	for i := 0; i < t.MaxRetries; i++ {
		req.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		resp, err = t.Transport.RoundTrip(req)
		if err != nil {
			return nil, err
		}

		if !t.shouldRetry(resp.StatusCode) {
			break
		}

		zap.L().Info(fmt.Sprintf("Retrying: %s %s %s", req.Method, req.URL, resp.Status))
		time.Sleep(time.Second * time.Duration(i+1))
	}

	return resp, err
}

func (t *Transport) shouldRetry(statusCode int) bool {
	for _, code := range t.RetryCodes {
		if statusCode == code {
			return true
		}
	}
	return false
}
