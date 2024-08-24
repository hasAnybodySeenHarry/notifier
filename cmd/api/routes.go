package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)
	router.NotFound = http.HandlerFunc(app.notFound)

	// meta
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheck)

	// notifications
	router.HandlerFunc(http.MethodGet, "/v1/notifications/subscribe", app.mustAuth(app.notificationSubscriberHandler))

	return app.enableCORS(router)
}
