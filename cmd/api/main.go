package main

import (
	"log"
	"os"
	"sync"

	"harry2an.com/notifier/internal/data"
	"harry2an.com/notifier/internal/metrics"
	"harry2an.com/notifier/internal/redis"
	"harry2an.com/notifier/internal/rpc"
)

type application struct {
	config  config
	wg      sync.WaitGroup
	logger  *log.Logger
	server  *server
	models  *data.Models
	grpc    *rpc.Clients
	metrics *metrics.Metrics
	clients *redis.Clients
}

func main() {
	var cfg config
	loadConfig(&cfg)

	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, conn, consumers, client, err := initDependencies(cfg, l)
	if err != nil {
		l.Fatalln(err)
	}
	defer func() {
		err = closeClient(db)
		if err != nil {
			l.Fatalln(err)
		}

		err = conn.Close()
		if err != nil {
			l.Fatalln(err)
		}

		err = consumers.Close()
		if err != nil {
			l.Fatalln(err)
		}

		err = client.Close()
		if err != nil {
			l.Fatalln(err)
		}
	}()

	server := &server{
		users: make(map[int64]*user, 0),
	}

	app := application{
		config:  cfg,
		logger:  l,
		server:  server,
		models:  data.New(db, cfg.db.database),
		grpc:    rpc.NewClients(conn),
		metrics: metrics.Register(),
		clients: redis.New(client),
	}

	var servers sync.WaitGroup
	servers.Add(3)

	go func() {
		defer servers.Done()
		if err := app.consume(consumers); err != nil {
			app.logger.Fatalln("Topic consuming stopped with error:", err)
		}
	}()

	go func() {
		defer servers.Done()
		if err := app.serve(); err != nil {
			app.logger.Fatalln("HTTP server stopped with error:", err)
		}
	}()

	go func() {
		defer servers.Done()
		if err := app.subscribeAndListen(); err != nil {
			app.logger.Fatalln("Redis Pub/Sub stopped with error:", err)
		}
	}()

	servers.Wait()
	app.logger.Println("All services have stopped gracefully.")
}
