package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type config struct {
	port     int
	env      string
	pub      publisher
	db       db
	grpcAddr string
}

type publisher struct {
	host string
	port int
}

type db struct {
	username string
	password string
	database string
	host     string
	port     int
}

func loadConfig(cfg *config) {
	flag.StringVar(&cfg.grpcAddr, "grpcAddr", os.Getenv("GRPC-ADDR"), "The address of the gRPC server")

	flag.IntVar(&cfg.port, "port", getEnvInt("PORT", 4000), "The port that the server listens at")
	flag.StringVar(&cfg.env, "env", os.Getenv("ENV"), "The environment of the server")

	flag.StringVar(&cfg.pub.host, "pub-host", os.Getenv("PUB_HOST"), "The address to connect to Kafka node")
	flag.IntVar(&cfg.pub.port, "pub-port", getEnvInt("PUB_PORT", 9092), "The port to connect to Kafka node")

	flag.StringVar(&cfg.db.username, "mongo-username", os.Getenv("MONGO_USERNAME"), "The username to connect to the mongo database")
	flag.StringVar(&cfg.db.password, "mongo-password", os.Getenv("MONGO_PASSWORD"), "The password to connect to the mongo database")
	flag.StringVar(&cfg.db.database, "mongo-database", os.Getenv("MONGO_DATABASE"), "The mongo database to connect to")
	flag.StringVar(&cfg.db.host, "mongo-host", os.Getenv("MONGO_HOST"), "The address to connect to the mongo database")
	flag.IntVar(&cfg.db.port, "mongo-port", getEnvInt("MONGO_PORT", 27017), "The port to connect to the mongo database")

	flag.Parse()
}

func getEnvInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		fmt.Printf("Invalid value for environment variable %s: %s\n", key, valueStr)
		return defaultValue
	}

	return value
}
