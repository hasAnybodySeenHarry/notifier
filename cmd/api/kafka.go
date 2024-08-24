package main

import (
	"log"
	"time"

	"github.com/IBM/sarama"
)

func openConsumerGroup(brokers []string, groupID string, retries int) (sarama.ConsumerGroup, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Retry.Backoff = 2
	config.Consumer.Group.Rebalance.GroupStrategies = append(config.Consumer.Group.Rebalance.GroupStrategies, sarama.NewBalanceStrategyRoundRobin())

	var err error
	var consumerGroup sarama.ConsumerGroup

	for i := 0; i < retries; i++ {
		consumerGroup, err = sarama.NewConsumerGroup(brokers, groupID, config)
		if err == nil {
			return consumerGroup, nil
		}

		log.Printf("Failed to connect to Kafka: %s. Retrying in 5 seconds...", err)
		time.Sleep(5 * time.Second)
	}

	return nil, err
}
