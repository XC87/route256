package domain

type ListItem struct {
	Sku_id  int64
	Count   uint64
	Product Product
}

type Item struct {
	Sku_id int64
	Count  uint64
}

type ItemsMap map[int64]Item
