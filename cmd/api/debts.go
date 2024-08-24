package main

import (
	"harry2an.com/notifier/internal/consumer"
)

func (app *application) notifyDebt(lenderID int64, borrowerID int64, message []byte, d *consumer.Debt) {
	notiID, err := app.models.Debts.Insert(d)
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
