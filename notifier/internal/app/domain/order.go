package domain

type OrderInfo struct {
	OrderId   int64  `json:"order_id" valid:"type(int64),required"`
	NewStatus string `json:"new_status" valid:"type(string),required"`
}
