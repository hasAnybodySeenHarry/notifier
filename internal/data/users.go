package data

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserStateModel struct {
	db *mongo.Database
}

type UserState struct {
	UserID         int64              `bson:"user_id"`
	LastSentNotiID primitive.ObjectID `bson:"last_sent_noti_id"`
}

func (m *UserStateModel) UpSert(u *UserState) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{"user_id": u.UserID}
	update := bson.M{"$set": bson.M{"last_sent_noti_id": u.LastSentNotiID}}
	opts := options.Update().SetUpsert(true)

	res, err := m.db.Collection(userStateCol).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return primitive.NilObjectID, err
	}

	if res.UpsertedID == nil {
		return primitive.NilObjectID, nil
	}

	upsertedID, ok := res.UpsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, nil
	}

	return upsertedID, nil
}

func (m *UserStateModel) GetLastInsertedID(userID int64) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	var res UserState

	err := m.db.Collection(userStateCol).FindOne(ctx, filter).Decode(&res)
	if err != nil {
		switch {
		case errors.Is(err, mongo.ErrNoDocuments):
			return primitive.NilObjectID, ErrNoRecords
		default:
			return primitive.NilObjectID, err
		}
	}

	return res.LastSentNotiID, nil
}
