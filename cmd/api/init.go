package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func initDependencies(
	cfg config, logger *log.Logger,
) (
	db *mongo.Client,
	conn *grpc.ClientConn,
	consumers sarama.ConsumerGroup,
	client *redis.Client,
	err error,
) {
	var clientErr, grpcErr, consumerErr, redisErr error

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		consumers, consumerErr = openConsumerGroup([]string{fmt.Sprintf("%s:%d", cfg.pub.host, cfg.pub.port)}, "default-consumer-group", 5)
		if consumerErr != nil {
			logger.Printf("Failed to join the consumers: %v", consumerErr)
		} else {
			logger.Println("Successfully joined the consumers")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		db, clientErr = openDB(cfg.db, 6)
		if clientErr != nil {
			logger.Printf("Failed to connect to the database: %v", clientErr)
		} else {
			logger.Println("Successfully connected to the database")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		conn, grpcErr = openGRPC(cfg.grpcAddr)
		if grpcErr != nil {
			log.Printf("Failed to connect to the gRPC server: %v", grpcErr)
		} else {
			logger.Println("Successfully connected to the gRPC server")
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		client, redisErr = openRedis(&cfg.redis, 2, 6)
		if redisErr != nil {
			log.Printf("Failed to connect to the redis database: %v", redisErr)
		} else {
			logger.Println("Successfully connected to the redis database")
		}
	}()

	wg.Wait()

	if clientErr != nil {
		return nil, nil, nil, nil, clientErr
	}
	if grpcErr != nil {
		return nil, nil, nil, nil, grpcErr
	}
	if consumerErr != nil {
		return nil, nil, nil, nil, consumerErr
	}
	if redisErr != nil {
		return nil, nil, nil, nil, redisErr
	}

	return db, conn, consumers, client, nil
}
