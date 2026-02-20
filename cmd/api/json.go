package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// struct validator
var Validate *validator.Validate

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

func writeJson(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {
	// allowing body to have <=1mb data
	maxByte := 1_048_578 // bytes of 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxByte))

	// decoding data
	decoder := json.NewDecoder(r.Body)
	// disallowing unknown feilds
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)

}

func writeJSONEroor(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJson(w, status, &envelope{Error: message})
}
