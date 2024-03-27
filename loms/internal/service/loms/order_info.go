package order_usecase

import (
	"context"
	"route256.ozon.ru/project/loms/internal/model"
)

func (s *Service) OrderInfo(ctx context.Context, id int64) (*model.Order, error) {
	if id <= 0 {
		return nil, ErrOrderInvalid
	}

	order, err := s.OrderRepository.OrderInfo(ctx, id)
	if err != nil {
		return nil, ErrOrderNotFound
	}

	return order, nil
}
