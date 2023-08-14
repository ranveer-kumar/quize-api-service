package handlers

import (
	"encoding/json"
	// "log"
	"net/http"
	"quize-api-service/internal/db"
	log "quize-api-service/internal/logger"
	"quize-api-service/internal/models"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetQuizeHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	quize, err := db.GetQuizeByID(id)
	if err != nil {
		// logger.Printf("Error fetching quize: %v", err)
		log.Error("Error fetching quize: ", err)
		http.Error(w, "Error fetching quize", http.StatusInternalServerError)
		return
	}
	if quize == nil {
		log.Info("Quize not found with id: " + id)
		http.Error(w, "Quize not found", http.StatusNotFound)
		return
	}

	// Convert the time fields to IST
	quize.ConvertTimeToIST()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(quize)
}

func CreateQuizeHandler(w http.ResponseWriter, r *http.Request) {
	var newQuize models.Quize
	err := json.NewDecoder(r.Body).Decode(&newQuize)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	newQuize.ID = primitive.NewObjectID().Hex() // Generate a new unique ID

	currentTime := time.Now()
	newQuize.CreatedAt = currentTime
	newQuize.UpdatedAt = currentTime
	// newQuize.CreatedAt = time.Now().UTC()

	insertedID, err := db.InsertQuize(&newQuize)
	if err != nil {
		http.Error(w, "Error creating quize", http.StatusInternalServerError)
		return
	}

	// Respond with the inserted ID
	response := map[string]string{"insertedID": insertedID}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func DeleteQuizeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := db.DeleteQuize(id)
	if err != nil {
		// log.Printf("Error deleting quiz: %v", err)
		log.Error("Error deleting quiz: ", err)
		http.Error(w, "Error deleting quiz", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func UpdateQuizeHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"] // Extract the quize ID from the URL parameters

	// Find the existing quize document
	existingQuize, err := db.GetQuizeByID(id)
	if err != nil {
		// log.Printf("Error fetching existing quize: %v", err)
		log.Error("Error fetching existing quize: ", err)
		http.Error(w, "Error fetching existing quize", http.StatusInternalServerError)
		return
	}
	if existingQuize == nil {
		http.Error(w, "Quize not found", http.StatusNotFound)
		return
	}

	var updatedQuize models.Quize
	err = json.NewDecoder(r.Body).Decode(&updatedQuize)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Set additional fields for the updated quiz
	updatedQuize.ID = id
	updatedQuize.UpdatedAt = time.Now()
	updatedQuize.CreatedAt = existingQuize.CreatedAt

	// Update the quiz in the database
	err = db.UpdateQuize(id, &updatedQuize)
	if err != nil {
		// log.Printf("Error updating quize: %v", err)
		log.Error("Error updating quize: ", err)
		http.Error(w, "Error updating quize", http.StatusInternalServerError)
		return
	}

	// Respond with a success message
	response := map[string]string{"message": "Quiz updated successfully"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func ListQuizeHandler(w http.ResponseWriter, r *http.Request) {

	// vars := mux.Vars(r)

	// // Get query parameters for pagination
	// pageSizeParam := vars["pageSize"]
	// pageNumberParam := vars["pageNumber"]
	pageSizeParam := r.URL.Query().Get("pageSize")
	pageNumberParam := r.URL.Query().Get("pageNumber")

	// Convert query parameters to integers
	pageSize, err := strconv.Atoi(pageSizeParam)
	if err != nil || pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	pageNumber, err := strconv.Atoi(pageNumberParam)
	if err != nil || pageNumber <= 0 {
		pageNumber = 1 // Default page number
	}

	// Fetch the paginated list of quizzes from the database
	quizes, err := db.ListQuizes(pageSize, pageNumber)
	if err != nil {
		// log.Printf("Error fetching quizzes: %v", err)
		log.Error("Error fetching quizzes: ", err)
		http.Error(w, "Error fetching quizzes", http.StatusInternalServerError)
		return
	}

	// Get the total count of quizzes for pagination calculations
	totalCount, err := db.GetTotalQuizeCount() // Implement this function in db package
	if err != nil {
		// log.Printf("Error getting total quiz count: %v", err)
		log.Error("Error getting total quiz count: ", err)
		http.Error(w, "Error getting total quiz count", http.StatusInternalServerError)
		return
	}

	// Calculate total number of pages
	totalPages := (totalCount + pageSize - 1) / pageSize

	// Create pagination metadata
	pagination := map[string]interface{}{
		"totalCount": totalCount,
		"pageSize":   pageSize,
		"pageNumber": pageNumber,
		"totalPages": totalPages,
	}

	// Prepare the response payload
	responsePayload := map[string]interface{}{
		"pagination": pagination,
		"quizes":     quizes,
	}

	// Respond with the paginated list of quizzes
	w.Header().Set("Content-Type", "application/json")
	encoder := json.NewEncoder(w)
	encoder.Encode(responsePayload)
}
