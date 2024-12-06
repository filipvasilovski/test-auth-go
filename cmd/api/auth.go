package main

import (
	"net/http"
	"os"
)

type LoginPayload struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,max=4,no_uppercase,no_numbers,no_symbols"`
	AccessToken string `json:"access_token"`
}

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var payload LoginPayload

	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err := Validate.Struct(payload)

	if err != nil {
		errors := GetValidatorErrorMessages(err)
		app.jsonResponse(w, http.StatusBadRequest, errors)
		return

	}

	pwd := os.Getenv("PASSWORD")
	email := os.Getenv("EMAIL")

	if payload.Email != email {
		// TODO: Make re-usable
		errPayload := []map[string]string{{"Param": "Email", "Message": "Wrong email address"}}
		app.jsonResponse(w, http.StatusUnauthorized, errPayload)
		return
	}

	if payload.Password != pwd {
		errPayload := []map[string]string{{"Param": "Password", "Message": "Wrong password"}}
		app.jsonResponse(w, http.StatusUnauthorized, errPayload)
		return
	}

	tokenString, err := createToken(payload.Email)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return

	}

	payload.AccessToken = tokenString

	if err := app.jsonResponse(w, http.StatusOK, payload); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
