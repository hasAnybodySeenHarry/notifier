package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"harry2an.com/notifier/internal/consumer"
)

type DebtModel struct {
	collection *mongo.Collection
}

func (m *DebtModel) Insert(d *consumer.Debt) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	document := bson.M{
		"id":         d.Data.ID,
		"lender":     d.Data.Lender,
		"borrower":   d.Data.Borrower,
		"category":   d.Data.Category,
		"total":      d.Data.Total,
		"created_at": d.Data.CreatedAt,
	}

	res, err := m.collection.InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}
