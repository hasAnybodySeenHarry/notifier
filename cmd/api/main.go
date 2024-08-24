package main

import (
	"log"
	"os"
	"sync"

	"harry2an.com/notifier/internal/data"
)

type application struct {
	config  config
	wg      sync.WaitGroup
	logger  *log.Logger
	server  *server
	models  *data.Models
	clients data.Clients
}

func main() {
	var cfg config
	loadConfig(&cfg)

	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	dbClient, conn, err := initDependencies(cfg, l)
	if err != nil {
		l.Fatalln(err)
	}
	defer func() {
		err = closeClient(dbClient)
		if err != nil {
			l.Fatalln(err)
		}
	}()

	server := &server{
		users: make(map[int64]*client, 0),
	}

	app := application{
		config:  cfg,
		logger:  l,
		server:  server,
		models:  data.New(dbClient, cfg.db.database),
		clients: data.NewClients(conn),
	}

	var servers sync.WaitGroup
	servers.Add(2)

	go func() {
		defer servers.Done()
		if err := app.consume(); err != nil {
			app.logger.Fatalln("Topic consuming stopped with error:", err)
		}
	}()

	go func() {
		defer servers.Done()
		if err := app.serve(); err != nil {
			app.logger.Fatalln("HTTP server stopped with error:", err)
		}
	}()

	servers.Wait()
	app.logger.Println("All services have stopped gracefully.")
}
