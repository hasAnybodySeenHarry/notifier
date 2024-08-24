package consumer

import (
	"time"

	"github.com/google/uuid"
	"harry2an.com/notifier/internal/core"
)

const (
	TransactionCreated = "transaction_created"
)

type Transaction struct {
	Metadata metadata        `json:"metadata"`
	Data     transactionData `json:"data"`
}

type transactionData struct {
	ID          int64       `json:"id"`
	Lender      core.Entity `json:"lender"`
	Borrower    core.Entity `json:"borrower"`
	DebtID      int64       `json:"debt_id"`
	Amount      float64     `json:"amount"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"created_at"`
	Version     uuid.UUID   `json:"version"`
}
