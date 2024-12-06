package main

import (
	"net/http"
	"os"
)

func (app *application) getCodeHandler(w http.ResponseWriter, r *http.Request) {
	c := os.Getenv("TC_CODE")

	if err := app.jsonResponse(w, http.StatusOK, c); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
