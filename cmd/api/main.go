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
	clients *clients
	models  *data.Models
}

func main() {
	var cfg config
	loadConfig(&cfg)

	l := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	dbClient, err := initDependencies(cfg, l)
	if err != nil {
		l.Fatalln(err)
	}
	defer func() {
		err = closeClient(dbClient)
		if err != nil {
			l.Fatalln(err)
		}
	}()

	clients := &clients{
		users:   make(map[int64]*client, 0),
		dummyID: 1,
	}

	app := application{
		config:  cfg,
		logger:  l,
		clients: clients,
		models:  data.New(dbClient, cfg.db.database),
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
