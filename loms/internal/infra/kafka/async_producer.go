package async_producer

import (
	"context"
	"fmt"
	"github.com/dnwe/otelsarama"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
	"route256.ozon.ru/pkg/logger"
	"route256.ozon.ru/project/loms/internal/model"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type AsyncProducer struct {
	producer sarama.AsyncProducer
}

func (p *AsyncProducer) ProduceMessage(ctx context.Context, message *model.KafkaMessage) error {
	msg := &sarama.ProducerMessage{
		Topic: message.Destination,
		Key:   sarama.StringEncoder(message.Key),
		Value: sarama.ByteEncoder(message.Data),
		Headers: []sarama.RecordHeader{
			{
				Key:   []byte("app-name"),
				Value: []byte("loms"),
			},
		},
		Timestamp: time.Now(),
	}
	otel.GetTextMapPropagator().Inject(ctx, otelsarama.NewProducerMessageCarrier(msg))

	select {
	case <-ctx.Done():
		return ctx.Err()
	case p.producer.Input() <- msg:
		return nil
	}
}

func NewKafkaAsyncProducer(ctx context.Context, wg *sync.WaitGroup, brokers []string, opts ...Option) (*AsyncProducer, error) {
	config := PrepareConfig(opts...)
	asyncProducer, err := sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("NewSyncProducer failed: %w", err)
	}

	asyncProducer = otelsarama.WrapAsyncProducer(config, asyncProducer)

	go func() {
		<-ctx.Done()
		err := asyncProducer.Close()
		if err != nil {
			fmt.Errorf("can't close kafka produced failed: %w", err)
		}
	}()

	producer := &AsyncProducer{
		producer: asyncProducer,
	}
	runKafkaSuccess(ctx, producer, wg)
	runKafkaErrors(ctx, producer, wg)

	return producer, nil
}

func runKafkaSuccess(ctx context.Context, asyncProducer *AsyncProducer, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		successCh := asyncProducer.producer.Successes()
		for {
			select {
			case <-ctx.Done():
				zap.L().Info("Kafka success ctx closed")
				return
			case msg := <-successCh:
				ctx = otel.GetTextMapPropagator().Extract(context.Background(), otelsarama.NewProducerMessageCarrier(msg))
				if msg == nil {
					zap.L().With(logger.GetTraceFields(ctx)...).Info("Kafka success chan closed")
					return
				}
				zap.L().With(logger.GetTraceFields(ctx)...).Sugar().Infof("Kafka success key: %q, partition: %d, offset: %d\n", msg.Key, msg.Partition, msg.Offset)
			}
		}
	}()
}

func runKafkaErrors(ctx context.Context, asyncProducer *AsyncProducer, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		errCh := asyncProducer.producer.Errors()

		for {
			select {
			case <-ctx.Done():
				zap.L().Info("Kafka error ctx closed")
				return
			case msgErr := <-errCh:
				if msgErr == nil {
					zap.L().Info("Kafka error chan closed")
					return
				}
				zap.S().Infof("Kafka error err %s, topic: %q, offset: %d\n", msgErr.Err, msgErr.Msg.Topic, msgErr.Msg.Offset)
			}
		}
	}()
}
