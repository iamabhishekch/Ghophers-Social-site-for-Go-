package main

import (
	"net/http"
)

// in this file we are handling error:
// for user which will go basic like Internal server error
// for developer, shows compelte error which helps to debugs

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err)

	writeJSONEroor(w, http.StatusInternalServerError, "the server encoutered a problem")
}

// bad request error

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONEroor(w, http.StatusBadRequest, err.Error())

}

// conflictResponse

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorf("conflict response", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONEroor(w, http.StatusBadRequest, err.Error())

}

// status not found

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONEroor(w, http.StatusNotFound, "not found")

}
