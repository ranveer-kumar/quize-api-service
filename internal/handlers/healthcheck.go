package handlers

import (
	"net/http"
	"quize-api-service/internal/logger"
)

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("health file pinged")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintln(w, "status: available")
// 	fmt.Fprintf(w, "environment: %s\n", app.config.env)
// 	fmt.Fprintf(w, "version: %s\n", version)
// }
