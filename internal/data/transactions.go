package data

// type TransactionModel struct {
// 	db *mongo.Database
// }

// func (m *TransactionModel) Insert(t *consumer.Transaction) (primitive.ObjectID, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// 	defer cancel()

// 	document := Transaction{
// 		ID:          primitive.NewObjectID(),
// 		Lender:      t.Data.Lender,
// 		Borrower:    t.Data.Borrower,
// 		DebtID:      t.Data.DebtID,
// 		Description: t.Data.Description,
// 		Amount:      t.Data.Amount,
// 		CreatedAt:   t.Data.CreatedAt,
// 	}

// 	res, err := m.collection.InsertOne(ctx, document)
// 	if err != nil {
// 		return primitive.NilObjectID, err
// 	}

// 	return res.InsertedID.(primitive.ObjectID), nil
// }
