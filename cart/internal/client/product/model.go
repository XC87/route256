package product

type ProductGetProductRequest struct {
	// Токен для доступа, нужно получить в Телеграмме у преподавателя
	Token *string `json:"token,omitempty"`
	// Уникальный id товара
	Sku *int64 `json:"sku,omitempty"`
}

type ProductGetProductResponse struct {
	Name  *string `json:"name,omitempty"`
	Price *uint32 `json:"price,omitempty"`
	//	Price *int64  `json:"price,omitempty"`
}

type ProductListSkusRequest struct {
	// Токен для доступа, нужно получить в Телеграмме у преподавателя
	Token *string `json:"token,omitempty"`
	// Начиная с какой sku выводить список (сама sku не включается в список)
	StartAfterSku *int64 `json:"startAfterSku,omitempty"`
	// Количество sku, которые надо вернуть
	Count *int64 `json:"count,omitempty"`
}

type ProductListSkusResponse struct {
	Skus []int64 `json:"skus,omitempty"`
}
