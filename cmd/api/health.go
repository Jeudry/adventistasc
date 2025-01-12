package main

import (
	"net/http"
)

// @Summary		Health Check
// @Description	Returns the health status of the application
// @Tags			health
// @Accept			json
// @Produce		json
// @Success		200	{object}	map[string]string	"Health status, environment, and version"
// @Failure		500	{object}	error				"Internal server error"
// @Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "Ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		writeJsonError(w, http.StatusInternalServerError, err.Error())
	}
}
