package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/asaskevich/govalidator"
	"github.com/dnwe/otelsarama"
	"go.opentelemetry.io/otel"
	semconv "go.opentelemetry.io/otel/semconv/v1.9.0"
	"go.opentelemetry.io/otel/trace"
	"route256.ozon.ru/project/notifier/internal/app/domain"
)

var _ sarama.ConsumerGroupHandler = (*orderChangeHandler)(nil)

type ConsumerHandler interface {
	Handle(data domain.OrderInfo) error
}

type orderChangeHandler struct {
	handler ConsumerHandler
}

func newOrderChangedHandler(handler ConsumerHandler) sarama.ConsumerGroupHandler {
	return otelsarama.WrapConsumerGroupHandler(&orderChangeHandler{handler: handler})
}

func (h *orderChangeHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			ctx := otel.GetTextMapPropagator().Extract(context.Background(), otelsarama.NewConsumerMessageCarrier(msg))
			tr := otel.Tracer("consumer")
			_, span := tr.Start(ctx, "consume message", trace.WithAttributes(
				semconv.MessagingOperationProcess,
			))
			defer span.End()

			var orderInfo domain.OrderInfo
			err := json.Unmarshal(msg.Value, &orderInfo)
			if err != nil {
				fmt.Errorf("failed unmarshal kafka orderInfo: %w", err)
				return nil
			}

			if _, err = govalidator.ValidateStruct(orderInfo); err != nil {
				fmt.Errorf("received invalid msg: %w", err)
				return nil
			}

			err = h.handler.Handle(orderInfo)
			if err != nil {
				fmt.Errorf("error from orderChangeHandler: %w", err)
				return nil
			}
			session.MarkMessage(msg, "")
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
