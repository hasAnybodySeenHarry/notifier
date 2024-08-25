package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"harry2an.com/notifier/internal/consumer"
	"harry2an.com/notifier/internal/core"
)

type TransactionModel struct {
	collection *mongo.Collection
}

type Transaction struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Lender      core.Entity        `bson:"lender"`
	Borrower    core.Entity        `bson:"borrower"`
	DebtID      int64              `bson:"debt_id"`
	Description string             `bson:"desc"`
	Amount      float64            `bson:"amount"`
	CreatedAt   time.Time          `bson:"created_at"`
}

func (m *TransactionModel) Insert(t *consumer.Transaction) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	document := Transaction{
		ID:          primitive.NewObjectID(),
		Lender:      t.Data.Lender,
		Borrower:    t.Data.Borrower,
		DebtID:      t.Data.DebtID,
		Description: t.Data.Description,
		Amount:      t.Data.Amount,
		CreatedAt:   t.Data.CreatedAt,
	}

	res, err := m.collection.InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}
