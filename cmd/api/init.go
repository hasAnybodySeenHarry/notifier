package main

import (
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc"
)

func initDependencies(cfg config, logger *log.Logger) (*mongo.Client, *grpc.ClientConn, error) {
	var clientErr, grpcErr error
	var client *mongo.Client
	var conn *grpc.ClientConn

	var wg sync.WaitGroup

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
		return nil, nil, clientErr
	}

	return client, conn, nil
}
