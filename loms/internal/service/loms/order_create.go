package order_usecase

import (
	"context"

	"route256.ozon.ru/project/loms/internal/model"
)

func (s *Service) OrderCreate(ctx context.Context, order *model.Order) (int64, error) {
	if order == nil || len(order.Items) == 0 {
		return 0, ErrOrderInvalid
	}

	if err := s.checkItemsExists(ctx, order.Items); err != nil {
		return 0, ErrOrderSkuNotFound
	}

	orderId, err := s.OrderRepository.OrderCreate(ctx, order)
	if err != nil {
		return 0, err
	}

	order.Id = orderId
	if err = s.checkAndReserveStock(ctx, order); err != nil {
		return 0, ErrOrderCantReserve
	}

	order.ChangeStatus(model.AwaitingPayment)
	err = s.OrderRepository.OrderUpdate(ctx, order)
	if err != nil {
		return 0, err
	}

	return orderId, nil
}

func (s *Service) checkItemsExists(ctx context.Context, items []model.Item) error {
	for _, item := range items {
		_, err := s.StockRepository.GetCountBySku(ctx, item.SKU)
		if err != nil {
			return ErrOrderSkuNotFound
		}
	}
	return nil
}

func (s *Service) checkAndReserveStock(ctx context.Context, order *model.Order) error {
	err := checkIfCanReserve(s.StockRepository, order.Items)
	if err != nil {
		order.ChangeStatus(model.Failed)
		err = s.OrderRepository.OrderUpdate(ctx, order)
		if err != nil {
			return err
		}
		return ErrOrderCantReserve
	}

	err = s.StockRepository.Reserve(ctx, order.Items)
	if err != nil {
		return err
	}

	return nil
}
