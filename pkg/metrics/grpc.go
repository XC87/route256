package metrics

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"strconv"
	"time"
)

var (
	grpcRequestCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "grpc_requests_total",
		Help: "Total number of gRPC requests",
	})

	grpcRequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "grpc_request_duration_seconds",
		Help:    "Duration of gRPC requests in seconds",
		Buckets: prometheus.DefBuckets,
	}, []string{"method", "status"})
)

func UnaryServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	start := time.Now()
	grpcRequestCounter.Inc()
	resp, err := handler(ctx, req)
	duration := time.Since(start)
	status, _ := status.FromError(err)
	statusCode := runtime.HTTPStatusFromCode(status.Code())
	grpcRequestDuration.WithLabelValues(info.FullMethod, strconv.Itoa(statusCode)).Observe(duration.Seconds())

	return resp, err
}

func UnaryClientInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	start := time.Now()
	grpcRequestCounter.Inc()
	err := invoker(ctx, method, req, reply, cc, opts...)
	status, _ := status.FromError(err)
	statusCode := runtime.HTTPStatusFromCode(status.Code())
	grpcRequestDuration.WithLabelValues(method, strconv.Itoa(statusCode)).Observe(time.Since(start).Seconds())

	return err
}
