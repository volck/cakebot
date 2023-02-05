package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	//healthcheck
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	//individual date
	router.HandlerFunc(http.MethodGet, "/v1/cake", app.getCakeHandler)
	router.HandlerFunc(http.MethodGet, "/v1/cake/:date", app.getDateHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/cake/:date", app.updateCakesHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/cake/:date", app.unsetCakeHandler)

	//plural
	router.HandlerFunc(http.MethodGet, "/v1/cakes", app.listCakesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/chef/:uid", app.getChef)

	router.HandlerFunc(http.MethodPost, "/v1/cakes", app.createCakeHandler)

	//PATCH

	router.HandlerFunc(http.MethodPost, "/v1/cake", app.insertCakeHandler)

	// USERS
	router.HandlerFunc(http.MethodGet, "/v1/users", app.getAllUsersHandler)

	// Return the httprouter instance.
	return router
}
