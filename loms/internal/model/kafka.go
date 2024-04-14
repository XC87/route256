package model

type KafkaMessage struct {
	Key         string
	Destination string
	Data        []byte
}

type OrderStatusInfo struct {
	OrderId   int64       `json:"order_id"`
	NewStatus OrderStatus `json:"new_status"`
}
