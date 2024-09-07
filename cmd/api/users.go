package main

import (
	"context"
	"time"

	"harry2an.com/notifier/cmd/proto/users"
)

func (app *application) getUserForToken(token string) (*users.GetUserResponse, error) {
	req := &users.GetUserRequest{Token: token}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := app.grpc.Users.GetUserForToken(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
