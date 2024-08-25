package consumer

import (
	"time"

	"github.com/google/uuid"
	"harry2an.com/notifier/internal/core"
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

func (t *Transaction) TransactionToMap() map[string]interface{} {
	return map[string]interface{}{
		"metadata": t.Metadata.MetadataToMap(),
		"data": map[string]interface{}{
			"id":          t.Data.ID,
			"lender":      core.EntityToMap(&t.Data.Lender),
			"borrower":    core.EntityToMap(&t.Data.Borrower),
			"debt_id":     t.Data.DebtID,
			"amount":      t.Data.Amount,
			"description": t.Data.Description,
			"created_at":  t.Data.CreatedAt.Format(time.RFC3339),
			"version":     t.Data.Version.String(),
		},
	}
}
