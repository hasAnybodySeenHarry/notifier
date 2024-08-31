package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
)

func (app *application) consume(cg sarama.ConsumerGroup) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	relay := make(chan os.Signal, 1)
	signal.Notify(relay, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		for {
			if err := cg.Consume(ctx, []string{"debts", "transactions"}, app); err != nil {
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
		log.Println(string(message.Value))
		app.processMessage(message)
		session.MarkMessage(message, "")
	}
	return nil
}
