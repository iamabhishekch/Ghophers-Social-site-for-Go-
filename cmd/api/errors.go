package main

import (
	"log"
	"net/http"
)

// in this file we are handling error:
// for user which will go basic like Internal server error
// for developer, shows compelte error which helps to debugs

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s, path: %s, error %s", r.Method, r.URL.Path, err)
	writeJSONEroor(w, http.StatusInternalServerError, "the server encoutered a problem")
}

// bad request error

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request response: %s, path: %s, error %s", r.Method, r.URL.Path, err)
	writeJSONEroor(w, http.StatusBadRequest, err.Error())

}


// status not found

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request response: %s, path: %s, error %s", r.Method, r.URL.Path, err)
	writeJSONEroor(w, http.StatusNotFound, "not found")

}
