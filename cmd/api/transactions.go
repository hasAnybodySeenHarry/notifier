package main

import (
	"harry2an.com/notifier/internal/consumer"
	"harry2an.com/notifier/internal/data"
)

func (app *application) notifyTransaction(lenderID int64, borrowerID int64, t *consumer.Transaction) {
	n := &data.Notification{
		Type:    consumer.TransactionCreated,
		Payload: t.TransactionToMap(),
	}

	notiID, err := app.models.Notifications.Insert(n)
	if err != nil {
		app.logger.Printf("Error encountered while inserting %v", err)
	}

	err = app.broadcastNotiIDForUsers([]int64{borrowerID, lenderID}, notiID, consumer.TransactionCreated)
	if err != nil {
		app.logger.Printf("Failed to broadcast notifications to the Pub/Sub: %v", err)
	}
}
