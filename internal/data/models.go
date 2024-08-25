package data

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrNoRecords = errors.New("error document not found")
)

type Models struct {
	Debts        DebtModel
	Transactions TransactionModel
	UserStates   UserStateModel
}

func New(c *mongo.Client, dbName string) *Models {
	db := c.Database(dbName)

	return &Models{
		Debts: DebtModel{
			collection: db.Collection("notification_collections"),
		},
		Transactions: TransactionModel{
			collection: db.Collection("notification_collections"),
		},
		UserStates: UserStateModel{
			collection: db.Collection("userstate_collections"),
		},
	}
}

type Notification struct {
	ID          primitive.ObjectID     `bson:"_id" json:"id"`
	OtherFields map[string]interface{} `bson:",inline"`
}

func (m *Models) GetDifferedNotifications(notiID primitive.ObjectID, userID int64) ([]*Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": bson.M{"$gt": notiID},
		"$or": []bson.M{
			{"lender.id": userID},
			{"borrower.id": userID},
		},
	}

	cursor, err := m.Debts.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var notifications []*Notification

	for cursor.Next(ctx) {
		var n Notification
		if err := cursor.Decode(&n); err != nil {
			return nil, err
		}
		notifications = append(notifications, &n)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}
