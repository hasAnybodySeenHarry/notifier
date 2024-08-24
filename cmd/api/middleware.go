package main

import (
	"net/http"
	"strings"
)

const AuthorizationHeader = "Authorization"

func (app *application) mustAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", AuthorizationHeader)

		h := r.Header.Get(AuthorizationHeader)
		if h == "" {
			app.invalidAuthToken(w, r)
			return
		}

		segs := strings.Split(h, " ")
		if len(segs) != 2 || segs[0] != "Bearer" {
			app.invalidAuthToken(w, r)
			return
		}

		token := segs[1]
		user, err := app.getUserForToken(token)
		if err != nil {
			app.handleGRPCError(w, r, err)
			return
		}

		next.ServeHTTP(w, app.setUser(r, user))
	}
}

func (app *application) enableCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Expected-Version")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		h.ServeHTTP(w, r)
	})
}
