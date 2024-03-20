package handlers

type ItemAddRequest struct {
	UserId int64  `json:"user_id" valid:"type(int64),required"`
	SkuId  int64  `json:"sku_id" valid:"type(int64),required"`
	Count  uint64 `json:"count" valid:"type(uint64),required"`
}

type ItemDeleteRequest struct {
	UserId int64 `json:"user_id" valid:"type(int64),required"`
	SkuId  int64 `json:"sku_id" valid:"type(int64),required"`
}

type CartGetRequest struct {
	UserId int64 `json:"user_id" valid:"type(int64),required"`
}

type OrderCheckoutRequest struct {
	UserId int64 `json:"user_id" valid:"type(int64),required"`
}
type OrderCheckoutResponse struct {
	OrderId int64 `json:"orderID" valid:"type(int64),required"`
}

type CartClearRequest struct {
	UserId int64 `json:"user_id" valid:"type(int64),required"`
}
