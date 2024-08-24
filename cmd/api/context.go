package main

import (
	"context"
	"net/http"

	"harry2an.com/notifier/cmd/proto/users"
)

type contextKey string

const userCtx = contextKey("user")

func (app *application) setUser(r *http.Request, u *users.GetUserResponse) *http.Request {
	ctx := context.WithValue(r.Context(), userCtx, u)
	return r.WithContext(ctx)
}

func (app *application) getUser(r *http.Request) *users.GetUserResponse {
	u, ok := r.Context().Value(userCtx).(*users.GetUserResponse)
	if !ok {
		panic("invalid user pointer inside the context")
	}

	return u
}
