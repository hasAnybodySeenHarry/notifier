package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type clients struct {
	users   map[int64]*client
	dummyID int64
	mu      sync.Mutex
}

type client struct {
	latest primitive.ObjectID
	conn   *websocket.Conn
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (app *application) notificationSubscriberHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		app.logger.Printf("Error encountered while upgrading as the web socket: %v", err)
		return
	}
	defer conn.Close()

	userID := app.clients.dummyID
	app.clients.dummyID++

	app.addWebSocketUser(userID, conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	app.clients.mu.Lock()
	defer app.clients.mu.Unlock()

	if err := app.removeWebSocketUser(userID); err != nil {
		app.logger.Printf("Error encountered while removing the client with ID %d %v", userID, err)
	}
}

func (app *application) addWebSocketUser(userID int64, conn *websocket.Conn) {
	app.clients.mu.Lock()
	defer app.clients.mu.Unlock()

	if _, exists := app.clients.users[userID]; !exists {
		app.clients.users[userID] = &client{
			latest: primitive.NilObjectID,
			conn:   conn,
		} // query the last acessed noti id from database
	}
}

func (app *application) removeWebSocketUser(userID int64) error {
	app.clients.mu.Lock()
	defer app.clients.mu.Unlock()

	client, ok := app.clients.users[userID]
	if !ok {
		return fmt.Errorf("user with ID %d got already removed", userID)
	}

	log.Printf("saving the latest sent event id %d", client.latest)

	delete(app.clients.users, userID)
	return client.conn.Close()
}

func (app *application) sendMessageToClient(clientID int64, notiID primitive.ObjectID, message []byte) error {
	app.clients.mu.Lock()
	client, ok := app.clients.users[clientID]
	app.clients.mu.Unlock()

	if !ok {
		return nil
	}

	if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
		app.logger.Printf("Failed to send message to user with ID: %d: %v", clientID, err)
		if removeErr := app.removeWebSocketUser(clientID); removeErr != nil {
			app.logger.Printf("Failed to remove user with ID: %d: %v", clientID, removeErr)
		}
		return err
	}

	client.latest = notiID

	return nil
}
