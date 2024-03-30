package loms

import (
	"context"
	"errors"
	order_usecase "route256.ozon.ru/project/loms/internal/service/loms"

	"route256.ozon.ru/project/loms/internal/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	servicepb "route256.ozon.ru/project/loms/pkg/api/v1"
)

func (s *Server) OrderCreate(ctx context.Context, request *servicepb.OrderCreateRequest) (*servicepb.OrderCreateResponse, error) {
	id, err := s.impl.OrderCreate(ctx, repackOrder(request))
	if err != nil {
		switch {
		case errors.Is(err, order_usecase.ErrOrderCantReserve):
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		case errors.Is(err, order_usecase.ErrOrderSkuNotFound):
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		default:
			return nil, status.Errorf(codes.Internal, "Failed to create order: %v", err)
		}
	}
	return &servicepb.OrderCreateResponse{OrderId: id}, nil
}

func repackOrder(in *servicepb.OrderCreateRequest) *model.Order {
	items := make([]model.Item, 0)
	for _, item := range in.Items {
		items = append(items, model.Item{
			SKU:   item.Sku,
			Count: item.Count,
		})
	}
	return &model.Order{
		Status: model.New,
		User:   in.User,
		Items:  items,
	}
}
