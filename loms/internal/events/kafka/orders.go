package kafka_events

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"log"
	"route256.ozon.ru/project/loms/internal/config"
	kafka "route256.ozon.ru/project/loms/internal/infra/kafka"
	"route256.ozon.ru/project/loms/internal/model"
	order_usecase "route256.ozon.ru/project/loms/internal/service/loms"
	"strconv"
	"sync"
	"time"
)

func RegisterEvents(ctx context.Context, wg *sync.WaitGroup, config *config.Config, eventManager order_usecase.EventManagers) {
	eventProducer := createProducer(ctx, wg, config)

	eventManager.Subscribe("order-events", func(ctx context.Context, data any) error {
		message, err := buildStatusMessage(config.LomsKafkaOrderStatusTopic, data)
		if err != nil || message == nil {
			return err
		}
		err = eventProducer.ProduceMessage(ctx, message)
		if err != nil {
			return fmt.Errorf("failed to emit order status changed event: %w", err)
		}
		return nil
	})
}

func createProducer(ctx context.Context, wg *sync.WaitGroup, config *config.Config) *kafka.AsyncProducer {
	eventProducer, err := kafka.NewKafkaAsyncProducer(
		ctx,
		wg,
		config.LomsKafkaBrokers,
		kafka.WithRequiredAcks(sarama.WaitForAll),
		kafka.WithMaxRetries(5),
		kafka.WithRetryBackoff(10*time.Millisecond),
		kafka.WithProducerFlushMessages(3),
		kafka.WithProducerFlushFrequency(5*time.Second),
	)
	if err != nil {
		log.Fatalf("cant start event producer: %v", err)
	}

	return eventProducer
}

func buildStatusMessage(OrderStatusChangedTopicName string, data any) (*model.KafkaMessage, error) {
	order, errb := data.(*model.Order)
	if errb == false {
		// если не сконвертилось то скорей всего чет левое, тогда пропускаем
		return nil, nil
	}
	orderInfo := model.OrderStatusInfo{
		OrderId:   order.Id,
		NewStatus: order.Status,
	}

	bytes, err := json.Marshal(orderInfo)
	if err != nil {
		return nil, fmt.Errorf("failed to build kafka status message: %w", err)
	}
	message := &model.KafkaMessage{
		Key:         strconv.FormatInt(order.Id, 10),
		Destination: OrderStatusChangedTopicName,
		Data:        bytes,
	}

	return message, nil
}
