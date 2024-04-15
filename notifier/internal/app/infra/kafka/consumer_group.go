package kafka

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/IBM/sarama"
)

type OrderChangedHandler struct {
	consumerGroup sarama.ConsumerGroup
	handler       sarama.ConsumerGroupHandler
	topics        []string
}

type Config struct {
	Brokers      []string
	TopicNames   []string
	ConsumerName string
}

func NewOrderChangeHandler(
	config Config,
	hanlder ConsumerHandler,
) (*OrderChangedHandler, error) {
	handler := newOrderChangedHandler(hanlder)

	cg, err := newConsumerGroup(
		config.Brokers,
		config.ConsumerName,
	)
	if err != nil {
		log.Fatal(err)
	}

	return &OrderChangedHandler{
		consumerGroup: cg,
		handler:       handler,
		topics:        config.TopicNames,
	}, nil
}

func (c *OrderChangedHandler) Run(ctx context.Context, wg *sync.WaitGroup) {
	c.runCGErrorHandler(ctx, wg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("[consumer-group] run")

		for {
			if err := c.consumerGroup.Consume(ctx, c.topics, c.handler); err != nil {
				log.Printf("Error from consume: %v\n", err)
			}
			if ctx.Err() != nil {
				c.consumerGroup.Close()
				log.Printf("[consumer-group]: ctx closed: %s\n", ctx.Err().Error())
				return
			}
		}
	}()
}

func (c *OrderChangedHandler) runCGErrorHandler(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case chErr, ok := <-c.consumerGroup.Errors():
				if !ok {
					log.Println("[cg-error] error: chan closed")
					return
				}

				log.Printf("[cg-error] error: %s\n", chErr)
			case <-ctx.Done():
				log.Printf("[cg-error]: ctx closed: %s\n", ctx.Err().Error())
				return
			}
		}
	}()
}

func newConsumerGroup(brokers []string, groupID string, opts ...Option) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Version = sarama.MaxVersion
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.ResetInvalidOffsets = true
	config.Consumer.Group.Heartbeat.Interval = 3 * time.Second
	config.Consumer.Group.Session.Timeout = 60 * time.Second
	config.Consumer.Group.Rebalance.Timeout = 60 * time.Second
	//
	config.Consumer.Return.Errors = true

	config.Consumer.Offsets.AutoCommit.Enable = false
	config.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second

	config.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{sarama.NewBalanceStrategyRoundRobin()}

	for _, opt := range opts {
		opt.Apply(config)
	}

	cg, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return cg, nil
}
