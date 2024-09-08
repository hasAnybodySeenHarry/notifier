package consumer

import (
	"time"

	"github.com/google/uuid"
	"harry2an.com/notifier/internal/core"
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
	Version   uuid.UUID   `json:"version"`
}

func (d *Debt) DebtToMap() map[string]interface{} {
	return map[string]interface{}{
		"metadata": d.Metadata.MetadataToMap(),
		"data": map[string]interface{}{
			"id":         d.Data.ID,
			"lender":     core.EntityToMap(&d.Data.Lender),
			"borrower":   core.EntityToMap(&d.Data.Borrower),
			"category":   d.Data.Category,
			"total":      d.Data.Total,
			"created_at": d.Data.CreatedAt.Format(time.RFC3339),
			"version":    d.Data.Version.String(),
		},
	}
}
