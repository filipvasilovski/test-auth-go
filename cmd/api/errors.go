package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"

	"github.com/go-playground/validator/v10"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) unauthorizedResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusUnauthorized, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflict error: %s path: %s error: %s", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusConflict, err.Error())
}

type ApiError struct {
	Param   string
	Message string
}

func msgForTag(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return "Invalid email address"
	case "no_uppercase":
		return "Uppercase letters are not allowed!"
	case "no_numbers":
		return "Numbers are not allowed!"
	case "no_symbols":
		return "Symbols are not allowed!"
	case "max":
		return fmt.Sprintf("%s cannot exceed %s characters", fe.Field(), fe.Param())
	}
	return fe.Error()
}

var lowercaseAndNumberRegex = regexp.MustCompile("^[a-z0-9]+$")

var (
	noUppercaseRegex = regexp.MustCompile("^[^A-Z]+$")      // No uppercase letters
	noNumbersRegex   = regexp.MustCompile("^[^0-9]+$")      // No numbers
	noSymbolsRegex   = regexp.MustCompile("^[a-zA-Z0-9]+$") // No special symbols
)

// CustomNoUppercaseValidator checks if a string contains no uppercase letters.
func CustomNoUppercaseValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return noUppercaseRegex.MatchString(value)
}

// CustomNoNumbersValidator checks if a string contains no numbers.
func CustomNoNumbersValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return noNumbersRegex.MatchString(value)
}

// CustomNoSymbolsValidator checks if a string contains no symbols (only letters and numbers allowed).
func CustomNoSymbolsValidator(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return noSymbolsRegex.MatchString(value)
}
func GetValidatorErrorMessages(e error) []ApiError {
	var ve validator.ValidationErrors

	if errors.As(e, &ve) {
		out := make([]ApiError, len(ve))
		for i, fe := range ve {
			out[i] = ApiError{fe.Field(), msgForTag(fe)}
		}

		return out
	}

	return nil
}
