package data

import "go.mongodb.org/mongo-driver/mongo"

type Models struct {
	Debts        DebtModel
	Transactions TransactionModel
}

func New(c *mongo.Client, dbName string) *Models {
	db := c.Database(dbName)

	return &Models{
		Debts: DebtModel{
			collection: db.Collection("debts_collections"),
		},
		Transactions: TransactionModel{
			collection: db.Collection("transactions_collections"),
		},
	}
}
