package model

type OrderStatus string

const (
	New             OrderStatus = "new"
	AwaitingPayment OrderStatus = "awaiting payment"
	Failed          OrderStatus = "failed"
	Paid            OrderStatus = "payed"
	Cancelled       OrderStatus = "cancelled"
)

type Item struct {
	SKU   uint32
	Count uint64
}

type Order struct {
	Id     int64
	Status OrderStatus
	User   int64
	Items  []Item
}

func (o *Order) ChangeStatus(status OrderStatus) {
	o.Status = status
}
