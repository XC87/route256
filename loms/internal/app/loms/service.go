package loms

import (
	"context"
	"route256.ozon.ru/project/loms/internal/model"

	_ "route256.ozon.ru/project/loms/internal/model"
	servicepb "route256.ozon.ru/project/loms/pkg/api/v1"
)

var _ servicepb.LomsServer = (*Server)(nil)

type LomsService interface {
	OrderCreate(ctx context.Context, order *model.Order) (int64, error)
	OrderInfo(ctx context.Context, id int64) (*model.Order, error)
	OrderPay(ctx context.Context, id int64) error
	OrderCancel(ctx context.Context, id int64) error
	StockInfo(ctx context.Context, sku uint32) (uint64, error)
}

type Server struct {
	servicepb.UnimplementedLomsServer
	impl LomsService
}

func NewServer(impl LomsService) *Server {
	return &Server{impl: impl}
}

var orderStatusMap = map[model.OrderStatus]servicepb.OrderInfoResponse_StatusEnum{
	model.New:             servicepb.OrderInfoResponse_new,
	model.AwaitingPayment: servicepb.OrderInfoResponse_awaiting_payment,
	model.Failed:          servicepb.OrderInfoResponse_failed,
	model.Paid:            servicepb.OrderInfoResponse_paid,
	model.Cancelled:       servicepb.OrderInfoResponse_cancelled,
}

func mapStatus(orderStatus model.OrderStatus) servicepb.OrderInfoResponse_StatusEnum {
	return orderStatusMap[orderStatus]
}
