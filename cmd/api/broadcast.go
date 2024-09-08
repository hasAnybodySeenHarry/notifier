package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"harry2an.com/notifier/internal/data"
	"harry2an.com/notifier/internal/redis"
)

func (app *application) broadcastNotiIDForUsers(users []int64, notiID primitive.ObjectID, notiType string) error {
	n := &redis.NotiBroadcast{
		ID:       notiID,
		Type:     notiType,
		ForUsers: users,
	}

	return app.clients.Users.Publish(n)
}

func (app *application) subscribeAndListen() error {
	relay := make(chan os.Signal, 1)
	signal.Notify(relay, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan struct{}, 1)
	var err error

	go func() {
		sub := app.clients.Users.Subscribe("noti-broadcast")
		defer func() {
			err = sub.Close()
		}()

		ch := sub.Channel()

		for {
			select {
			case <-ctx.Done():
				stop <- struct{}{}
				return
			case m := <-ch:
				var b redis.NotiBroadcast

				err := json.Unmarshal([]byte(m.Payload), &b)
				if err != nil {
					app.logger.Printf("Failed to decode payload %s as a NotiBroadcast", m.Payload)
				}

				active := make([]int64, 0, 2)
				for _, id := range b.ForUsers {
					_, ok := app.server.users[id]
					if ok {
						active = append(active, id)
					}
				}

				if len(active) == 0 {
					continue
				}

				n, err := app.models.Notifications.GetNotificationByID(b.ID)
				if err != nil {
					switch {
					case errors.Is(err, data.ErrNoRecords):
						app.logger.Printf("WARNING: notification with id %s is not found", b.ID.Hex())
					default:
						app.logger.Printf("Error retrieving notification with id %s", b.ID.Hex())
					}
				}

				message, err := json.Marshal(n.Payload)
				if err != nil {
					app.logger.Printf("WARNING: Failed to marshal notification: %v. Skipped sending it to clients.", err)
					continue
				}

				for _, id := range active {
					app.background(func() {
						sent, err := app.sendMessageToClient(id, n.ID, message)
						if err != nil {
							app.logger.Printf("%v", err)
						}

						if !sent {
							app.logger.Printf("WARNING: client with id %d was connected but is no longer being served by the server", id)
						}
					})
				}
			}
		}
	}()

	app.logger.Println("Subscribed and listening at the channel...")

	s := <-relay
	app.logger.Printf("Received signal %s and stopping redis pub/sub now", s.String())
	cancel()

	<-stop
	if err != nil {
		return err
	}

	app.logger.Println("Successfully closed the redis' pub/sub subscription.")
	return nil
}
