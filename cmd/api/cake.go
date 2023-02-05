package main

import (
	"cakebot/internal/data"
	"cakebot/internal/data/validator"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app application) insertCakeHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		UserID    string `json:"UserID"`
		When      string `json:"Cakeday"`
		Firstname string `json:"FIRSTNAME"`
		Lastname  string `json:"LASTNAME"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	cake := data.Cake{}
	cake.User_ID = input.UserID
	cake.When = input.When

	err = app.models.Cake.Insert(&cake)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	sendcake, err := app.models.Cake.GetByDate(cake.When)

	err = app.writeJSON(w, http.StatusOK, envelope{"cake": sendcake}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getCakeHandler(w http.ResponseWriter, r *http.Request) {

	cakes, err := app.models.Cake.GetCurrent()

	err = app.writeJSON(w, http.StatusOK, envelope{"cake": cakes}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) getDateHandler(w http.ResponseWriter, r *http.Request) {

	theDate, err := app.readDateParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	v := validator.New()
	const layout = "01.02.2006"

	_, validTimestamp := time.Parse(layout, theDate)
	v.Check(validTimestamp == nil, "date", "timestamp not valid. Valid formats are: DD.MM.YYYY")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	cakes, err := app.models.Cake.GetByDate(theDate)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"cake": cakes}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) listCakesHandler(w http.ResponseWriter, r *http.Request) {
	cakes, err := app.models.Cake.GetAll()
	// Send a JSON response containing the movie data.
	err = app.writeJSON(w, http.StatusOK, envelope{"cake": cakes}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) createCakeHandler(w http.ResponseWriter, r *http.Request) {

	// Declare an anonymous struct to hold the information that we expect to be in the
	// HTTP request body (note that the field names and types in the struct are a subset
	// of the Movie struct that we created earlier). This struct will be our *target
	// decode destination*.

	var input struct {
		UserID    string `json:"UserID"`
		When      string `json:"CAKEDAY"`
		Firstname string `json:"FIRSTNAME"`
		Lastname  string `json:"LASTNAME"`
	}
	// Initialize a new json.Decoder instance which reads from the request body, and
	// then use the Decode() method to decode the body contents into the input struct.
	// Importantly, notice that when we call Decode() we pass a *pointer* to the input
	// struct as the target decode destination. If there was an error during decoding,
	// we also use our generic errorResponse() helper to send the client a 400 Bad
	// Request response containing the error message.
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	input.When = app.convertTime(input.When)
	// Initialize a new Validator instance.
	v := validator.New()
	// Use the Check() method to execute our validation checks. This will add the
	// provided key and error message to the errors map if the check does not evaluate
	// to true. For example, in the first line here we "check that the title is not
	// equal to the empty string". In the second, we "check that the length of the title
	// is less than or equal to 500 bytes" and so on.
	v.Check(input.UserID != "", "UserID", "must be provided")
	v.Check(len(input.UserID) <= 6, "UserID", "UserID cannot exceed a99999")
	v.Check(validator.Matches(input.UserID, validator.UidRX), "UserID", "Not a valid userid. Must match ^a|A[0-9]{5}")
	v.Check(input.Firstname != "", "Firstname", "Firstname must be provided")
	v.Check(input.Lastname != "", "Lastname", "Lastname must be provided")
	//v.Check(input.When != data.CakeTime(), "When", "When cannot be empty")

	// values in the input.Genres slice are unique.
	//v.Check(validator.Unique(input.Genres), "genres", "must not contain duplicate values")
	// Use the Valid() method to see if any of the checks failed. If they did, then use
	// the failedValidationResponse() helper to send a response to the client, passing
	// in the v.Errors map.
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) getChef(w http.ResponseWriter, r *http.Request) {
	id, err := app.readUIdParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	v := validator.New()
	v.Check(id != "", "id", "must be provided")
	v.Check(len(id) < 50, "id", "must not be longer than 50")

	cake, err := app.models.Cake.GetUid(id)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"cakes": cake}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	// Dump the contents of the input struct in a HTTP response.
}

func (app *application) updateCakesHandler(w http.ResponseWriter, r *http.Request) {

	theDate, err := app.readDateParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	cake, err := app.models.Cake.GetByDate(theDate)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	type CakeInput struct {
		User_ID string `db:"USER_ID"`
		When    string `db:"CAKEDAY"`
		NewDate string `json:"NewDate"`
	}
	myInput := CakeInput{}

	err = app.readJSON(w, r, &myInput)
	if err != nil {
		app.badRequestResponse(w, r, err)
	}

	//lots of checks about input here
	cake.User_ID = myInput.User_ID
	cake.When = theDate
	cake.NewDate = myInput.NewDate

	v := validator.New()
	v.Check(cake.User_ID != "", "USER_ID", "UserID must be set")
	v.Check(cake.NewDate != "", "NEWDATE", "NEWDATE must be set")
	v.Check(cake.When != "", "DATE", "Date must be set")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// Pass the updated cake record to our new Update() method.
	err = app.models.Cake.Update(cake)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Write the updated movie record in a JSON response.
	returnCake, err := app.models.Cake.GetByDate(cake.NewDate)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"cake": returnCake}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) unsetCakeHandler(w http.ResponseWriter, r *http.Request) {

	theDate, err := app.readDateParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
	}

	cake, err := app.models.Cake.GetByDate(theDate)

	err = app.models.Cake.Unset(theDate)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"cake": cake, "status": "deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
