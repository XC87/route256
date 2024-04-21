package loms

import (
	"context"
	"github.com/pkg/errors"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/propagators/jaeger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"route256.ozon.ru/pkg/metrics"
	"route256.ozon.ru/project/cart/internal/domain"
	"route256.ozon.ru/project/cart/internal/service"
	servicepb "route256.ozon.ru/project/cart/pkg/api/v1"
)

type lomsGrpcClient struct {
	grpcClient servicepb.LomsClient
}

func NewLomsGrpcClient(ctx context.Context, serviceHost string) (service.LomsService, error) {
	conn, err := grpc.DialContext(
		ctx,
		serviceHost,
		grpc.WithUnaryInterceptor(metrics.UnaryClientInterceptor),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(
			otelgrpc.NewClientHandler(
				otelgrpc.WithPropagators(jaeger.Jaeger{}),
			),
		),
	)

	if err != nil {
		return nil, err
	}

	grpcClient := servicepb.NewLomsClient(conn)

	go func() {
		<-ctx.Done()
		zap.L().Info("Shutting down grpc client")
		if err = conn.Close(); err != nil {
			zap.L().Info("Failed to shutdown grpc client: ", zap.Error(err))
		}
	}()

	zap.L().Info("Loms grpc: connected to " + serviceHost)

	return &lomsGrpcClient{
		grpcClient: grpcClient,
	}, nil
}

func (l *lomsGrpcClient) CreateOrder(ctx context.Context, userId int64, items []domain.Item) (int64, error) {
	request := &servicepb.OrderCreateRequest{
		User:  userId,
		Items: make([]*servicepb.OrderItemCreateRequest, len(items)),
	}

	for i, item := range items {
		request.Items[i] = &servicepb.OrderItemCreateRequest{
			Sku:   uint32(item.Sku_id),
			Count: item.Count,
		}
	}

	//ctx = metadata.AppendToOutgoingContext(ctx, "x-auth", "123")
	response, err := l.grpcClient.OrderCreate(ctx, request)
	if err != nil {
		return 0, errors.Wrap(err, "loms.OrderCreate")
	}

	return response.OrderId, nil
}

func (l *lomsGrpcClient) GetStockInfo(ctx context.Context, sku uint32) (uint64, error) {
	request := &servicepb.StockInfoRequest{Sku: sku}
	//ctx = metadata.AppendToOutgoingContext(ctx, "x-auth", "123")
	response, err := l.grpcClient.StockInfo(ctx, request)
	if err != nil {
		return 0, errors.Wrap(err, "loms.GetStockInfo")
	}

	return response.Count, nil
}
