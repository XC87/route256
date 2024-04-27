package loms

import (
	"context"
	"errors"
	"google.golang.org/protobuf/types/known/emptypb"
	order_usecase "route256.ozon.ru/project/loms/internal/service/loms"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	servicepb "route256.ozon.ru/project/loms/pkg/api/v1"
)

func (s *Server) OrderCancel(ctx context.Context, request *servicepb.OrderCancelRequest) (*emptypb.Empty, error) {
	err := s.impl.OrderCancel(ctx, request.OrderId, request.UserId)
	if err != nil {
		if errors.Is(err, order_usecase.ErrOrderAlreadyCanceled) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
