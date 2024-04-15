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
		err = s.OrderRepository.OrderPay(ctx, id)
		if err == nil {
			s.EventManager.Trigger(ctx, "order-events", order)
		}
		return err
	case model.Paid:
		return ErrOrderAlreadyPaid
	default:
		return ErrOrderCantBePaid
	}
}
