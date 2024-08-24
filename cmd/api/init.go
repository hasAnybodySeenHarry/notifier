package main

import (
	"log"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
)

func initDependencies(cfg config, logger *log.Logger) (*mongo.Client, error) {
	var clientErr error
	var client *mongo.Client

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

	wg.Wait()

	if clientErr != nil {
		return nil, clientErr
	}

	return client, nil
}
