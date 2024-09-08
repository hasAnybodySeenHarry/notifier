package consumer

const (
	DebtCreated        = "debt_created"
	TransactionCreated = "transaction_created"
)

type Event struct {
	Metadata metadata `json:"metadata"`
	// Data     map[string]interface{} `json:"data"`
}

type metadata struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Type      string `json:"type"`
	Source    string `json:"source"`
	Version   string `json:"version"`
}

func (m *metadata) MetadataToMap() map[string]interface{} {
	return map[string]interface{}{
		"id":        m.ID,
		"timestamp": m.Timestamp,
		"type":      m.Type,
		"source":    m.Source,
		"version":   m.Version,
	}
}
