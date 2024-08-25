package main

import (
	"harry2an.com/notifier/internal/consumer"
	"harry2an.com/notifier/internal/data"
)

func (app *application) notifyTransaction(lenderID int64, borrowerID int64, message []byte, t *consumer.Transaction) {
	n := &data.Notification{
		Type:    consumer.TransactionCreated,
		Payload: t.TransactionToMap(),
	}

	notiID, err := app.models.Notifications.Insert(n)
	if err != nil {
		app.logger.Printf("Error encountered while inserting %v", err)
	}

	err = app.sendMessageToClient(borrowerID, notiID, message)
	if err != nil {
		app.logger.Printf("%v", err)
	}

	err = app.sendMessageToClient(lenderID, notiID, message)
	if err != nil {
		app.logger.Printf("%v", err)
	}
}
