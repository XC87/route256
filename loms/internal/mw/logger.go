package mw

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"log"
	"net/http"
)

func Logger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	reqMsg := req.(proto.Message)
	rawReq, _ := protojson.Marshal(reqMsg)
	log.Printf("request: method: %v, req: %v\n", info.FullMethod, string(rawReq))

	if resp, err = handler(ctx, req); err != nil {
		log.Printf("response: method: %v, err: %v\n", info.FullMethod, err)
		return
	}

	respMsg := resp.(proto.Message)
	rawResp, _ := protojson.Marshal(respMsg)
	log.Printf("response: method: %v, resp: %v\n", info.FullMethod, string(rawResp))

	return
}

func WithHTTPLoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("request: " + r.Method + " " + r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
