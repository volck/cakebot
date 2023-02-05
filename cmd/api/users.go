package main

import (
	"net/http"
)

func (app application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {

	users, err := app.models.Users.GetAll()

	err = app.writeJSON(w, http.StatusOK, envelope{"Users": users}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
