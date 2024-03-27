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

func (s *Server) OrderPay(ctx context.Context, request *servicepb.OrderPayRequest) (*emptypb.Empty, error) {
	err := s.impl.OrderPay(ctx, request.OrderId)
	if err != nil {
		if errors.Is(err, order_usecase.ErrOrderAlreadyPaid) {
			return nil, status.Errorf(codes.FailedPrecondition, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}
