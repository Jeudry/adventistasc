package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("internal server error: %s, path: %s error %s", r.Method, r.URL.Path, err)

	writeJsonError(w, http.StatusInternalServerError, "The server encountered a problem and could not process your request.")
}

func (app *application) badRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("bad request error: %s, path: %s error %s", r.Method, r.URL.Path, err)

	writeJsonError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFound(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s, path: %s error %s", r.Method, r.URL.Path, err)

	writeJsonError(w, http.StatusNotFound, "The requested resource could not be found.")
}
