package data

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"harry2an.com/notifier/internal/consumer"
	"harry2an.com/notifier/internal/core"
)

type DebtModel struct {
	collection *mongo.Collection
}

type Debt struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Lender    core.Entity        `bson:"lender"`
	Borrower  core.Entity        `bson:"borrower"`
	Category  string             `bson:"category"`
	Total     float64            `bson:"total"`
	CreatedAt time.Time          `bson:"created_at"`
}

func (m *DebtModel) Insert(d *consumer.Debt) (primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	document := Debt{
		ID:        primitive.NewObjectID(),
		Lender:    d.Data.Lender,
		Borrower:  d.Data.Borrower,
		Category:  d.Data.Category,
		Total:     d.Data.Total,
		CreatedAt: d.Data.CreatedAt,
	}

	res, err := m.collection.InsertOne(ctx, document)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}
