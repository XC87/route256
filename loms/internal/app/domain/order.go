package domain

import "errors"

var (
	ErrOrderInvalidId     = errors.New("order_id is required")
	ErrOrderInvalidStatus = errors.New("new_status is required")
)

type OrderInfo struct {
	OrderId   int64  `json:"order_id" valid:"type(int64),required"`
	NewStatus string `json:"new_status" valid:"type(string),required"`
}

func (oi *OrderInfo) Validate() error {
	if oi.OrderId == 0 {
		return ErrOrderInvalidId
	}
	if oi.NewStatus == "" {
		return ErrOrderInvalidStatus
	}
	return nil
}
