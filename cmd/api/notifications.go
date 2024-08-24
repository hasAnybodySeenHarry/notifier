package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type server struct {
	users map[int64]*client
	mu    sync.Mutex
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

	userID := app.getUser(r).Id
	app.addWebSocketUser(userID, conn)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}

	app.server.mu.Lock()
	defer app.server.mu.Unlock()

	if err := app.removeWebSocketUser(userID); err != nil {
		app.logger.Printf("Error encountered while removing the client with ID %d %v", userID, err)
	}
}

func (app *application) addWebSocketUser(userID int64, conn *websocket.Conn) {
	app.server.mu.Lock()
	defer app.server.mu.Unlock()

	if _, exists := app.server.users[userID]; !exists {
		app.server.users[userID] = &client{
			latest: primitive.NilObjectID,
			conn:   conn,
		} // query the last acessed noti id from database
	}
}

func (app *application) removeWebSocketUser(userID int64) error {
	app.server.mu.Lock()
	defer app.server.mu.Unlock()

	client, ok := app.server.users[userID]
	if !ok {
		return fmt.Errorf("user with ID %d got already removed", userID)
	}

	log.Printf("saving the latest sent event id %d", client.latest)

	delete(app.server.users, userID)
	return client.conn.Close()
}

func (app *application) sendMessageToClient(clientID int64, notiID primitive.ObjectID, message []byte) error {
	app.server.mu.Lock()
	client, ok := app.server.users[clientID]
	app.server.mu.Unlock()

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
