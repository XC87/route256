package model

type ProductStock struct {
	SKU        uint32 `json:"sku"`
	TotalCount uint64 `json:"total_count"`
	Reserved   uint64 `json:"reserved"`
}
