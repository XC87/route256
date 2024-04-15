package order_usecase

import (
	"context"
	"route256.ozon.ru/project/loms/internal/model"
)

func (s *Service) OrderCancel(ctx context.Context, id int64) error {
	if id == 0 {
		return ErrOrderInvalid
	}

	order, err := s.OrderRepository.OrderInfo(ctx, id)
	if err != nil {
		return ErrOrderNotFound
	}
	if order.Status == model.Cancelled {
		return ErrOrderAlreadyCanceled
	}

	err = checkIfCanUnReserve(s.StockRepository, order.Items)
	if err != nil {
		return ErrOrderCantUnReserve
	}
	s.StockRepository.UnReserve(ctx, order.Items)

	err = s.OrderRepository.OrderCancel(ctx, id)
	if err == nil {
		s.EventManager.Trigger(ctx, "order-events", order)
	}

	return err
}
