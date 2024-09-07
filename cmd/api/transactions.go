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

	client1, err := app.sendMessageToClient(borrowerID, notiID, message)
	if err != nil {
		app.logger.Printf("%v", err)
	}

	client2, err := app.sendMessageToClient(lenderID, notiID, message)
	if err != nil {
		app.logger.Printf("%v", err)
	}

	if !(client1 && client2) {
		usersToNotify := make([]int64, 0, 2)

		if !client1 {
			app.logger.Printf("Borrower with ID %d is not being served by the current instance", borrowerID)
			usersToNotify = append(usersToNotify, borrowerID)
		}

		if !client2 {
			app.logger.Printf("Lender with ID %d is not being served by the current instance", borrowerID)
			usersToNotify = append(usersToNotify, lenderID)
		}

		err := app.broadcastNotiIDForUsers(usersToNotify, notiID, consumer.DebtCreated)
		if err != nil {
			app.logger.Printf("Failed to broadcast notifications to the Pub/Sub: %v", err)
		}
	}
}
