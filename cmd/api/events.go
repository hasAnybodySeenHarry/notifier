package main

import (
	"encoding/json"

	"github.com/IBM/sarama"
	"harry2an.com/notifier/internal/consumer"
)

func (app *application) processMessage(message *sarama.ConsumerMessage) {
	var event consumer.Event

	// for performance sake, we avoid using a stream decoder at this stage
	// and instead, opted for the direct unmarshaling.
	err := json.Unmarshal(message.Value, &event)
	if err != nil {
		app.logger.Printf("Failed to parse the event's type: %v", err)
		return
	}

	switch event.Metadata.Type {
	case consumer.DebtCreated:
		var d consumer.Debt
		err := app.readJSON(message.Value, &d)
		if err != nil {
			app.logger.Printf("Failed to parse debt created event: %v", err)
			return
		}
		app.background(func() {
			app.notifyDebt(d.Data.Lender.ID, d.Data.Borrower.ID, message.Value, &d)
		})

	case consumer.TransactionCreated:
		var t consumer.Transaction
		err := app.readJSON(message.Value, &t)
		if err != nil {
			app.logger.Printf("Failed to parse transaction created event: %v", err)
			return
		}
		app.background(func() {
			app.notifyTransaction(t.Data.Lender.ID, t.Data.Borrower.ID, message.Value, &t)
		})

	default:
		app.logger.Printf("Unsupported event type: %s", event.Metadata.Type)
	}
}
