package order_usecase

import (
	"context"
	"route256.ozon.ru/project/loms/internal/model"
)

func (s *Service) OrderPay(ctx context.Context, id int64) error {
	if id == 0 {
		return ErrOrderInvalid
	}

	order, err := s.OrderRepository.OrderInfo(ctx, id)
	if err != nil {
		return ErrOrderNotFound
	}

	switch order.Status {
	case model.AwaitingPayment:
		return s.OrderRepository.OrderPay(ctx, id)
	case model.Paid:
		return ErrOrderAlreadyPaid
	default:
		return ErrOrderCantBePaid
	}
}
