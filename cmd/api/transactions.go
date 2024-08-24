package main

import (
	"harry2an.com/notifier/internal/consumer"
)

func (app *application) notifyTransaction(lenderID int64, borrowerID int64, message []byte, t *consumer.Transaction) {
	notiID, err := app.models.Transactions.Insert(t)
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
