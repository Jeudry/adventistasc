package main

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"log"
	"net/http"
)

func init() {
	Validate = validator.New(validator.WithRequiredStructEnabled())
}

var Validate *validator.Validate

func writeJson(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func writeJsonError(w http.ResponseWriter, status int, message string) {
	type envelope struct {
		Error string `json:"error"`
	}

	err := writeJson(w, status, envelope{
		Error: message,
	})

	if err != nil {
		log.Printf("error writing json error: %s", err)
	}

	return
}

func (app *application) jsonResponse(w http.ResponseWriter, status int, data any) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return writeJson(w, status, envelope{
		Data: data,
	})
}
