package loms

import (
	"context"
	servicepb "route256.ozon.ru/project/loms/pkg/api/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) StockInfo(ctx context.Context, request *servicepb.StockInfoRequest) (*servicepb.StockInfoResponse, error) {
	count, err := s.impl.StockInfo(ctx, request.Sku)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &servicepb.StockInfoResponse{Count: count}, nil
}
