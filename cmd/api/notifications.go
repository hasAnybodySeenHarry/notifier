package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"harry2an.com/notifier/internal/data"
)

type server struct {
	users map[int64]*client
	mu    sync.Mutex
}

type client struct {
	lastSent primitive.ObjectID
	conn     *websocket.Conn
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
	err = app.addWebSocketUser(userID, conn)
	if err != nil {
		app.logger.Printf("Error encountered while retrieving the last sent notification id with %v", err)
		return
	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			app.logger.Println("Error reading from the client or client just disconnected")
			break
		}
	}

	if err := app.removeWebSocketUser(userID); err != nil {
		app.logger.Printf("Error encountered while removing the client with ID %d %v", userID, err)
	}
}

func (app *application) addWebSocketUser(userID int64, conn *websocket.Conn) error {
	var notifications []*data.Notification

	lastID, err := app.models.UserStates.GetLastInsertedID(userID)
	if err != nil && !errors.Is(err, data.ErrNoRecords) {
		return err
	} else if err == nil {
		notifications, err = app.models.Notifications.GetNotifications(lastID, userID)
		if err != nil {
			return err
		}
	}

	// overriding the previous connection
	c := &client{
		lastSent: lastID,
		conn:     conn,
	}

	for _, notification := range notifications {
		message, err := json.Marshal(notification.Payload)
		if err != nil {
			app.logger.Printf("WARNING: Failed to marshal notification: %v", err)
			continue
		}

		err = conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			app.logger.Printf("WARNING: Failed to send message to userID %d: %v", userID, err)
		}

		c.lastSent = notification.ID
	}

	app.server.mu.Lock()
	defer app.server.mu.Unlock()
	app.server.users[userID] = c

	return nil
}

func (app *application) removeWebSocketUser(userID int64) error {
	app.server.mu.Lock()

	client, ok := app.server.users[userID]
	if !ok {
		return fmt.Errorf("user with ID %d got already removed", userID)
	}

	delete(app.server.users, userID)
	defer app.server.mu.Unlock()

	userState := &data.UserState{
		UserID:         userID,
		LastSentNotiID: client.lastSent,
	}

	// omitting the inserted id
	_, err := app.models.UserStates.UpSert(userState)
	if err != nil {
		app.logger.Printf("WARNING: Updating the last sent notification ID faild with %v", err)
	}

	err = client.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "You are being disconnected"))
	if err != nil {
		app.logger.Printf("WARNING: Failed to send close message to userID %d: %v", userID, err)
	}

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

	client.lastSent = notiID

	return nil
}
