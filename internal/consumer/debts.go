package consumer

import (
	"time"

	"github.com/google/uuid"
	"harry2an.com/notifier/internal/core"
)

const (
	DebtCreated = "debt_created"
)

type Debt struct {
	Metadata metadata `json:"metadata"`
	Data     debtData `json:"data"`
}

type debtData struct {
	ID        int64       `json:"id"`
	Lender    core.Entity `json:"lender"`
	Borrower  core.Entity `json:"borrower"`
	Category  string      `json:"category"`
	Total     float64     `json:"total"`
	CreatedAt time.Time   `json:"created_at"`
	Version   uuid.UUID   `json:"-"`
}
