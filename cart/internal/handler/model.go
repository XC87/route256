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
