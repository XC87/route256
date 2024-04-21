package mw

import (
	"go.uber.org/zap"
	"net/http"
	"route256.ozon.ru/pkg/logger"
)

type MiddlewareChain func(http.HandlerFunc) http.HandlerFunc

func BuildMiddleware(handler http.HandlerFunc, middlewareList []MiddlewareChain) http.HandlerFunc {
	if len(middlewareList) == 0 {
		return handler
	}
	wrapped := handler

	for _, middleware := range middlewareList {
		wrapped = middleware(wrapped)
	}

	return wrapped
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if (r.URL.Path == "/health") || (r.URL.Path == "/metrics") {
			next(w, r)
			return
		}
		zap.L().With(logger.GetTraceFields(r.Context())...).Info("request: " + r.Method + " " + r.URL.Path)
		next(w, r)
	}
}
