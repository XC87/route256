package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"route256.ozon.ru/project/notifier/internal/app/handlers"
	"route256.ozon.ru/project/notifier/internal/app/infra/kafka"
	"route256.ozon.ru/project/notifier/internal/config"
	"sync"
	"syscall"
)

func main() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

	config, err := config.GetConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	wg := &sync.WaitGroup{}

	kafkaConfig := kafka.Config{
		Brokers:      config.KafkaBrokers,
		TopicNames:   config.KafkaTopicNames,
		ConsumerName: config.KafkaConsumerName,
	}

	notifier := handlers.NewNotifier()
	kafkaEventListener, err := kafka.NewOrderChangeHandler(kafkaConfig, notifier)
	if err != nil {
		log.Fatal(err)
	}

	kafkaEventListener.Run(ctx, wg)

	<-ctx.Done()
	wg.Wait()
}
