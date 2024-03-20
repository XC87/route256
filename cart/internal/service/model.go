package service

type CartResponse struct {
	Items      []CartItem `json:"items"`
	TotalPrice uint32     `json:"total_price"`
}

type CartItem struct {
	SkuId int64  `json:"sku_id" valid:"type(int)"`
	Name  string `json:"name"`
	Count uint64 `json:"count"`
	Price uint32 `json:"price"`
}
