package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func openDB(cfg db, maxRetries int) (*mongo.Client, error) {
	var client *mongo.Client
	var err error

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d", cfg.username, cfg.password, cfg.host, cfg.port)
	clientOptions := options.Client().ApplyURI(uri)

	for i := 0; i < maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err = mongo.Connect(ctx, clientOptions)
		if err != nil {
			log.Printf("Failed to connect to the database: %v. Retrying in 5 seconds...", err)
			time.Sleep(5 * time.Second)
			continue
		}

		err = client.Ping(ctx, nil)
		if nil == err {
			return client, nil
		}

		log.Printf("Failed to ping the database: %v. Retrying in 3 seconds...", err)
		time.Sleep(3 * time.Second)
	}

	return client, nil
}

func closeClient(client *mongo.Client) error {
	if client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return client.Disconnect(ctx)
}
