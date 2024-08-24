package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"harry2an.com/notifier/internal/consumer"
)

type TransactionModel struct {
	collection *mongo.Collection
}

func (m *TransactionModel) Insert(t *consumer.Transaction) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	document := bson.M{
		"id":         t.Data.ID,
		"lender":     t.Data.Lender,
		"borrower":   t.Data.Borrower,
		"debt_id":    t.Data.DebtID,
		"desc":       t.Data.Description,
		"amount":     t.Data.Amount,
		"created_at": t.Data.CreatedAt,
	}

	res, err := m.collection.InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}
