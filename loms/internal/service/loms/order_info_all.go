package order_usecase

import (
	"context"
	"go.opentelemetry.io/otel"
	"google.golang.org/protobuf/types/known/emptypb"
	"route256.ozon.ru/project/loms/internal/model"
)

func (s *Service) OrderInfoAll(ctx context.Context, request *emptypb.Empty) ([]*model.Order, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "OrderInfoAll")
	defer span.End()

	order, err := s.OrderRepository.OrderInfoAll(ctx)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	return order, nil
}
