package routes

import (
	"quize-api-service/internal/handlers"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/v1/health", handlers.HealthCheckHandler).Methods("GET")
	router.HandleFunc("/v1/quizes/{id}", handlers.GetQuizeHandler).Methods("GET")
	router.HandleFunc("/v1/quizes", handlers.CreateQuizeHandler).Methods("POST")
	router.HandleFunc("/v1/quizes/{id}", handlers.UpdateQuizeHandler).Methods("PUT")
	router.HandleFunc("/v1/quizes", handlers.ListQuizeHandler).Methods("GET")
	router.HandleFunc("/v1/quizes/{id}", handlers.DeleteQuizeHandler).Methods("DELETE")

	// Define other routes for update, delete, etc.

	return router
}
