package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// router.HandlerFunc(http.MethodGet, "/v1/users", app.requireAdmin(app.listUserHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.requireOwnershipOrAdmin(app.getUserHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/users/:id", app.requireOwnershipOrAdmin(app.updateUserHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", app.requireOwnershipOrAdmin(app.deleteUserHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/users/:id/role", app.requireAdmin(app.updateUserRoleHandler))

	router.HandlerFunc(http.MethodPost, "/v1/users/activate", app.activateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/reset-password", app.updateUserPasswordHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/session", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/tokens/session", app.requireAuthenticatedUser(app.deleteAuthenticationTokenHandler))
	router.HandlerFunc(http.MethodPost, "/v1/tokens/activation", app.createActivationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/password-reset", app.createPasswordResetTokenHandler)

	router.Handler(http.MethodGet, "/debug/vars", app.requireAdmin(expvar.Handler().ServeHTTP))

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
