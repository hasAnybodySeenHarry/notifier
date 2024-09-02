package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed)
	router.NotFound = http.HandlerFunc(app.notFound)

	// meta
	router.Handler(http.MethodGet, "/metrics", promhttp.HandlerFor(app.metrics.Registry, promhttp.HandlerOpts{}))
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheck)

	// notifications
	router.HandlerFunc(http.MethodGet, "/v1/notifications/subscribe", app.mustAuth(app.notificationSubscriberHandler))

	return app.enableCORS(router)
}
