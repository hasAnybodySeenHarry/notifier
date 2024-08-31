package main

import (
	"fmt"
	"log"
	"sync"

	"github.com/IBM/sarama"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func initDependencies(cfg config, logger *log.Logger) (*mongo.Client, *grpc.ClientConn, sarama.ConsumerGroup, error) {
	var clientErr, grpcErr, consumerErr error

	var client *mongo.Client
	var conn *grpc.ClientConn
	var consumers sarama.ConsumerGroup

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
		client, clientErr = openDB(cfg.db, 6)
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

	wg.Wait()

	if clientErr != nil {
		return nil, nil, nil, clientErr
	}
	if grpcErr != nil {
		return nil, nil, nil, grpcErr
	}
	if consumerErr != nil {
		return nil, nil, nil, consumerErr
	}

	return client, conn, consumers, nil
}
