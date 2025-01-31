package main

import (
	"errors"
	"github.com/GoogleCloudPlatform/golang-samples/run/helloworld/internal/store"
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

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("not found error: %s, path: %s error %s", r.Method, r.URL.Path, err)

	writeJsonError(w, http.StatusNotFound, "The requested resource could not be found.")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("conflict error: %s, path: %s error %s", r.Method, r.URL.Path, err)

	writeJsonError(w, http.StatusConflict, "The requested resource could not be created due to a conflict.")
}

func (app *application) handleError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound):
		app.notFoundResponse(w, r, err)
	case errors.Is(err, store.ErrConflict):
		app.conflictResponse(w, r, err)
	default:
		app.internalServerError(w, r, err)
	}
}
