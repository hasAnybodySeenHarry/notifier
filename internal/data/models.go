package data

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	notificationCol = "notification_collection"
	userStateCol    = "userstate_collection"
)

var (
	ErrNoRecords = errors.New("error documents not found")
)

type Models struct {
	Notifications NotificationModel
	UserStates    UserStateModel
}

func New(c *mongo.Client, dbName string) *Models {
	db := c.Database(dbName)
	return &Models{
		Notifications: NotificationModel{db: db},
		UserStates:    UserStateModel{db: db},
	}
}
