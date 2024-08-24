package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

func InitDependencies(cfg config, logger *log.Logger) (sarama.ConsumerGroup, error) {
	// planned to extract consumer group init and mongo client and embed here

	var consumers sarama.ConsumerGroup
	var consumersErr error

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumers, consumersErr = openConsumerGroup([]string{fmt.Sprintf("%s:%d", cfg.pub.host, cfg.pub.port)}, "default-consumer-group", 5)
		if consumersErr != nil {
			logger.Printf("Failed to connect to the consumer groups: %v", consumersErr)
		} else {
			logger.Println("Successfully connected to the consumer groups")
		}
	}()

	wg.Wait()

	if consumersErr != nil {
		return nil, consumersErr
	}

	return consumers, nil
}

func (app *application) consume() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumerGroup, err := openConsumerGroup([]string{fmt.Sprintf("%s:%d", app.config.pub.host, app.config.pub.port)}, "default-consumer-group", 5)
	if err != nil {
		return err
	}
	defer consumerGroup.Close()
	app.logger.Println("Successfully connected to the consumer group")

	relay := make(chan os.Signal, 1)
	signal.Notify(relay, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, []string{"debts", "transactions"}, app); err != nil {
				app.logger.Printf("Error during consuming: %v", err)
				time.Sleep(2 * time.Second)
			}

			if err := ctx.Err(); err != nil {
				return
			}
		}
	}()

	<-relay
	cancel()
	app.logger.Println("Received signals to stop the consumer")

	app.logger.Println("Consumer group has been stopped")
	return nil
}

func (app *application) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (app *application) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (app *application) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		app.processMessage(message)
		session.MarkMessage(message, "")
	}
	return nil
}
