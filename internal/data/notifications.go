package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Notification struct {
	ID      primitive.ObjectID     `bson:"_id" json:"-"`
	Type    string                 `bson:"type" json:"type"`
	Payload map[string]interface{} `bson:"payload" json:"payload"`
}

type NotificationModel struct {
	db *mongo.Database
}

func (m *NotificationModel) GetNotifications(notiID primitive.ObjectID, userID int64) ([]*Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{
		"_id": bson.M{"$gt": notiID},
		"$or": []bson.M{
			{"payload.data.lender.id": userID},
			{"payload.data.borrower.id": userID},
		},
	}

	cursor, err := m.db.Collection(notificationCol).Find(ctx, filter)
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

func (m *NotificationModel) Insert(n *Notification) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	n.ID = primitive.NewObjectID()

	res, err := m.db.Collection(notificationCol).InsertOne(ctx, n)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (m *NotificationModel) GetNotificationByID(id primitive.ObjectID) (*Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var n Notification
	filter := bson.M{"_id": id}

	err := m.db.Collection(notificationCol).FindOne(ctx, filter).Decode(&n)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrNoRecords
		}
		return nil, err
	}

	return &n, nil
}

func (m *NotificationModel) GetLatestNotifications(userID int64) ([]*Notification, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"payload.data.lender.id": userID},
			{"payload.data.borrower.id": userID},
		},
	}

	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(5)

	cursor, err := m.db.Collection(notificationCol).Find(ctx, filter, opts)
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
