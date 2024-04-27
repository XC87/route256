package loms

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"route256.ozon.ru/project/loms/internal/model"
	servicepb "route256.ozon.ru/project/loms/pkg/api/v1"
)

func (s *Server) OrderInfo(ctx context.Context, request *servicepb.OrderInfoRequest) (*servicepb.OrderInfoResponse, error) {
	order, err := s.impl.OrderInfo(ctx, request.OrderId, request.UserId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return repackOrderToProto(order), nil
}

func repackOrderToProto(in *model.Order) *servicepb.OrderInfoResponse {
	orderItemInfoResponses := make([]*servicepb.OrderItemInfoResponse, 0, len(in.Items))
	for _, item := range in.Items {
		orderItemInfoResponses = append(orderItemInfoResponses, &servicepb.OrderItemInfoResponse{
			Sku:   item.SKU,
			Count: item.Count,
		})
	}
	return &servicepb.OrderInfoResponse{
		Id:     in.Id,
		Status: model.MapStatusToGrpc(in.Status),
		User:   in.User,
		Items:  orderItemInfoResponses,
	}
}
