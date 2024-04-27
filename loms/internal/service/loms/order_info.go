package order_usecase

import (
	"context"
	"go.opentelemetry.io/otel"
	"route256.ozon.ru/project/loms/internal/model"
)

func (s *Service) OrderInfo(ctx context.Context, id int64, userId int64) (*model.Order, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "OrderInfo")
	defer span.End()

	if id <= 0 {
		return nil, ErrOrderInvalid
	}

	order, err := s.OrderRepository.OrderInfo(ctx, id, userId)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	return order, nil
}
