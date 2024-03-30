package model

import (
	"route256.ozon.ru/project/loms/pkg/api/v1"
	"time"
)

type OrderStatus string

const (
	New             OrderStatus = "new"
	AwaitingPayment OrderStatus = "awaiting payment"
	Failed          OrderStatus = "failed"
	Paid            OrderStatus = "payed"
	Cancelled       OrderStatus = "cancelled"
)

type Item struct {
	SKU   uint32
	Count uint64
}

type Order struct {
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	Status    OrderStatus
	User      int64
	Items     []Item
}

func (o *Order) ChangeStatus(status OrderStatus) {
	o.Status = status
}

var orderStatusMap = map[OrderStatus]pb.OrderInfoResponse_StatusEnum{
	New:             pb.OrderInfoResponse_new,
	AwaitingPayment: pb.OrderInfoResponse_awaiting_payment,
	Failed:          pb.OrderInfoResponse_failed,
	Paid:            pb.OrderInfoResponse_paid,
	Cancelled:       pb.OrderInfoResponse_cancelled,
}

var orderStatusToID = map[OrderStatus]int32{
	New:             1,
	AwaitingPayment: 2,
	Failed:          3,
	Paid:            4,
	Cancelled:       5,
}

var iDToOrderStatus = map[int32]OrderStatus{
	1: New,
	2: AwaitingPayment,
	3: Failed,
	4: Paid,
	5: Cancelled,
}

func MapStatusToGrpc(orderStatus OrderStatus) pb.OrderInfoResponse_StatusEnum {
	return orderStatusMap[orderStatus]
}

func MapIdToStatus(orderStatus int32) OrderStatus {
	return iDToOrderStatus[orderStatus]
}

func MapStatusToId(orderStatus OrderStatus) int32 {
	return orderStatusToID[orderStatus]
}
