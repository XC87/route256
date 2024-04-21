package order_usecase

import (
	"context"
	"go.opentelemetry.io/otel"
	"route256.ozon.ru/project/loms/internal/model"
)

func (s *Service) OrderCreate(ctx context.Context, order *model.Order) (int64, error) {
	ctx, span := otel.Tracer("default").Start(ctx, "OrderCreate")
	defer span.End()

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
	s.EventManager.Trigger(ctx, "order-events", order)

	order.Id = orderId
	if err = s.checkAndReserveStock(ctx, order); err != nil {
		return 0, ErrOrderCantReserve
	}

	_, span = otel.Tracer("default").Start(ctx, "orderюChangeStatus")
	order.ChangeStatus(model.AwaitingPayment)
	err = s.OrderRepository.OrderUpdate(ctx, order)
	span.End()
	if err != nil {
		return 0, err
	}
	s.EventManager.Trigger(ctx, "order-events", order)

	return orderId, nil
}

func (s *Service) checkItemsExists(ctx context.Context, items []model.Item) error {
	ctx, span := otel.Tracer("default").Start(ctx, "order.checkItemsExists")
	defer span.End()

	for _, item := range items {
		_, err := s.StockRepository.GetCountBySku(ctx, item.SKU)
		if err != nil {
			return ErrOrderSkuNotFound
		}
	}
	return nil
}

func (s *Service) checkAndReserveStock(ctx context.Context, order *model.Order) error {
	ctx, span := otel.Tracer("default").Start(ctx, "order.checkAndReserveStock")
	defer span.End()

	err := checkIfCanReserve(s.StockRepository, order.Items)
	if err != nil {
		order.ChangeStatus(model.Failed)
		err = s.OrderRepository.OrderUpdate(ctx, order)
		if err != nil {
			return err
		}
		s.EventManager.Trigger(ctx, "order-events", order)
		return ErrOrderCantReserve
	}

	err = s.StockRepository.Reserve(ctx, order.Items)
	if err != nil {
		return err
	}

	return nil
}
