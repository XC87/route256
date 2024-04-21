package server

import (
	"context"
	"go.uber.org/zap"
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
		zap.L().Info("Shutting down http")
		if err := s.server.Shutdown(ctx); err != nil {
			zap.L().Info("Failed to shutdown http server: ", zap.Error(err))
		}
	}()

	zap.L().Info("Starting server")
	if err := s.server.ListenAndServe(); err != nil {
		zap.L().Fatal("error starting server", zap.Error(err))
	}
}
