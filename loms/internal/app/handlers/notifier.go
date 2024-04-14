package handlers

import (
	"log"
	"route256.ozon.ru/project/notifier/internal/app/domain"
)

type Notifier struct{}

func NewNotifier() *Notifier {
	return &Notifier{}
}

func (n *Notifier) Handle(data domain.OrderInfo) error {
	log.Printf("order id %d, has changed status to %s", data.OrderId, data.NewStatus)
	return nil
}
