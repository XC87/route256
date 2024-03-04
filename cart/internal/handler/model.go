package handler

type ItemAddRequest struct {
	UserId int64  `json:"user_id" valid:"type(int64),required"`
	SkuId  int64  `json:"sku_id" valid:"type(int64),required"`
	Count  uint16 `json:"count" valid:"type(uint16),required"`
}

type ItemDeleteRequest struct {
	UserId int64 `json:"user_id" valid:"type(int64),required"`
	SkuId  int64 `json:"sku_id" valid:"type(int64),required"`
}

type CartGetRequest struct {
	UserId int64 `json:"user_id" valid:"type(int64),required"`
}

type CartClearRequest struct {
	UserId int64 `json:"user_id" valid:"type(int64),required"`
}

type CartItem struct {
	SkuID int64  `json:"sku_id" valid:"type(int)"`
	Name  string `json:"name"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

type CartResponse struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}
