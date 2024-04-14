package server

import (
	"context"
	"log"
	"net/http"
	"route256.ozon.ru/project/cart/internal/config"
	"time"
)

type HTTPServer struct {
	server *http.Server
}

func NewHTTPServer(cartConfig *config.Config) *HTTPServer {
	return &HTTPServer{
		server: &http.Server{
			Addr:        cartConfig.CartHost,
			ReadTimeout: 10 * time.Second,
		},
	}
}

func (s *HTTPServer) Listen(ctx context.Context) {
	go func() {
		<-ctx.Done()
		log.Println("Shutting down http")
		if err := s.server.Shutdown(ctx); err != nil {
			log.Println("Failed to shutdown http server: ", err)
		}
	}()

	log.Println("Starting server")
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
