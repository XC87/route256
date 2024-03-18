package domain

type ListItem struct {
	Sku_id  int64
	Count   uint16
	Product Product
}

type Item struct {
	Sku_id int64
	Count  uint16
}

type ItemsMap map[int64]Item
