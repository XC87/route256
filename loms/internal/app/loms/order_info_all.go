package loms

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"route256.ozon.ru/project/loms/internal/model"
	servicepb "route256.ozon.ru/project/loms/pkg/api/v1"
)

func (s *Server) OrderInfoAll(ctx context.Context, request *emptypb.Empty) (*servicepb.OrderInfoAllResponse, error) {
	orderList, err := s.impl.OrderInfoAll(ctx, request)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	return repackOrderListToProto(orderList), nil
}

func repackOrderListToProto(list []*model.Order) *servicepb.OrderInfoAllResponse {
	orderItemInfoAllResponses := &servicepb.OrderInfoAllResponse{}
	for _, order := range list {
		orderItemInfoResponses := make([]*servicepb.OrderItemInfoResponse, 0, len(order.Items))
		for _, item := range order.Items {
			orderItemInfoResponses = append(orderItemInfoResponses, &servicepb.OrderItemInfoResponse{
				Sku:   item.SKU,
				Count: item.Count,
			})
		}
		orderItemInfoAllResponses.Items = append(orderItemInfoAllResponses.Items, &servicepb.OrderInfoResponse{
			Id:     order.Id,
			Status: model.MapStatusToGrpc(order.Status),
			User:   order.User,
			Items:  orderItemInfoResponses},
		)
	}

	return orderItemInfoAllResponses
}
