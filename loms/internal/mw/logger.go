package mw

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"net/http"
	"route256.ozon.ru/pkg/logger"
)

func Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	reqMsg := req.(proto.Message)
	rawReq, _ := protojson.Marshal(reqMsg)
	traceFields := logger.GetTraceFields(ctx)
	logger := zap.L().With(traceFields...).Sugar()

	logger.Infof("request: method: %v, req: %v", info.FullMethod, string(rawReq))

	if resp, err = handler(ctx, req); err != nil {
		logger.Errorf("response: method: %v, err: %v", info.FullMethod, err)
		return
	}

	respMsg := resp.(proto.Message)
	rawResp, _ := protojson.Marshal(respMsg)
	logger.Infof("response: method: %v, resp: %v", info.FullMethod, string(rawResp))

	return
}

func WithHTTPLoggingMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if (r.URL.Path == "/health") || (r.URL.Path == "/metrics") {
			next.ServeHTTP(w, r)
			return
		}
		zap.L().With(logger.GetTraceFields(r.Context())...).Info("request: " + r.Method + " " + r.URL.Path)
		next.ServeHTTP(w, r)
	}
}
