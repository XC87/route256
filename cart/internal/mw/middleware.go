package mw

import (
	"log"
	"net/http"
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
		log.Println("request: " + r.Method + " " + r.URL.Path)
		next(w, r)
	}
}
