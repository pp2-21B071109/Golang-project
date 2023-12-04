package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	// Use the requirePermission() middleware on each of the /v1/coins** endpoints,
	// passing in the required permission code as the first parameter.
	router.HandlerFunc(http.MethodGet, "/v1/coins", app.requirePermission("coins:read", app.listCoinsHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/coins/:id", app.requirePermission("coins:write", app.updateCoinHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/coins/:id", app.requirePermission("coins:write", app.deleteCoinHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	return app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router))))

}
