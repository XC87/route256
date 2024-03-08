package handler

import (
	"log"
	"net/http"
)

type middlewareChain func(http.HandlerFunc) http.HandlerFunc

func buildMiddleware(handler http.HandlerFunc, middlewareList []middlewareChain) http.HandlerFunc {
	if len(middlewareList) == 0 {
		return handler
	}
	wrapped := handler

	for _, middleware := range middlewareList {
		wrapped = middleware(wrapped)
	}

	return wrapped
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("request: " + r.Method + " " + r.URL.Path)
		next(w, r)
	}
}
