package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/asaskevich/govalidator"
	"route256.ozon.ru/project/notifier/internal/app/domain"

	"github.com/IBM/sarama"
)

var _ sarama.ConsumerGroupHandler = (*orderChangeHandler)(nil)

type ConsumerHandler interface {
	Handle(data domain.OrderInfo) error
}

type orderChangeHandler struct {
	handler ConsumerHandler
}

func newOrderChangedHandler(handler ConsumerHandler) *orderChangeHandler {
	return &orderChangeHandler{handler: handler}
}

func (h *orderChangeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			var orderInfo domain.OrderInfo
			err := json.Unmarshal(message.Value, &orderInfo)
			if err != nil {
				fmt.Errorf("failed unmarshal kafka orderInfo: %w", err)
				return nil
			}

			if _, err = govalidator.ValidateStruct(orderInfo); err != nil {
				fmt.Errorf("received invalid message: %w", err)
				return nil
			}

			err = h.handler.Handle(orderInfo)
			if err != nil {
				fmt.Errorf("error from orderChangeHandler: %w", err)
				return nil
			}
			session.MarkMessage(message, "")
			session.Commit()
		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *orderChangeHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *orderChangeHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}
