package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/faazabilamri7/mygram/database"
	"github.com/faazabilamri7/mygram/models"
	"github.com/gorilla/mux"
)

// CreateSocialMediaEntry handles the creation of a new social media entry
func CreateSocialMediaEntry(w http.ResponseWriter, r *http.Request) {
	var socialMedia models.SocialMedia
	err := json.NewDecoder(r.Body).Decode(&socialMedia)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Retrieve user ID from token or request context
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Set user ID for the social media entry
	socialMedia.UserID = userID

	// Save social media entry to the database
	err = database.GetDB().Create(&socialMedia).Error
	if err != nil {
		http.Error(w, "Failed to create social media entry", http.StatusInternalServerError)
		return
	}

	// Set appropriate response status and return the created social media entry
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(socialMedia)
}

// GetAllSocialMediaEntries handles fetching all social media entries from the logged-in user
func GetAllSocialMediaEntries(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from token or request context
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var socialMediaEntries []models.SocialMedia
	err = database.GetDB().Where("user_id = ?", userID).Find(&socialMediaEntries).Error
	if err != nil {
		http.Error(w, "Failed to fetch social media entries", http.StatusInternalServerError)
		return
	}

	// Set appropriate response status and return the fetched social media entries
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(socialMediaEntries)
}

// GetSocialMediaEntryByID handles fetching a specific social media entry by its ID from the logged-in user
func GetSocialMediaEntryByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from token or request context
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	socialMediaIDStr := mux.Vars(r)["socialMediaID"]
	socialMediaID, err := strconv.Atoi(socialMediaIDStr)
	if err != nil {
		http.Error(w, "Invalid social media ID", http.StatusBadRequest)
		return
	}

	var socialMediaEntry models.SocialMedia
	err = database.GetDB().Where("user_id = ? AND id = ?", userID, socialMediaID).First(&socialMediaEntry).Error
	if err != nil {
		http.Error(w, "Social media entry not found", http.StatusNotFound)
		return
	}

	// Set appropriate response status and return the fetched social media entry
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(socialMediaEntry)
}

// UpdateSocialMediaEntryByID handles updating a social media entry by its ID
func UpdateSocialMediaEntryByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from token or request context
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	socialMediaIDStr := mux.Vars(r)["socialMediaID"]
	socialMediaID, err := strconv.Atoi(socialMediaIDStr)
	if err != nil {
		http.Error(w, "Invalid social media ID", http.StatusBadRequest)
		return
	}

	var updatedSocialMedia models.SocialMedia
	err = json.NewDecoder(r.Body).Decode(&updatedSocialMedia)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the social media entry exists
	var existingSocialMedia models.SocialMedia
	err = database.GetDB().Where("user_id = ? AND id = ?", userID, socialMediaID).First(&existingSocialMedia).Error
	if err != nil {
		http.Error(w, "Social media entry not found", http.StatusNotFound)
		return
	}

	// Update social media entry fields
	existingSocialMedia.Name = updatedSocialMedia.Name
	existingSocialMedia.URL = updatedSocialMedia.URL

	// Save updated social media entry to the database
	err = database.GetDB().Save(&existingSocialMedia).Error
	if err != nil {
		http.Error(w, "Failed to update social media entry", http.StatusInternalServerError)
		return
	}

	// Set appropriate response status and return the updated social media entry
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingSocialMedia)
}

// DeleteSocialMediaEntryByID handles deleting a social media entry by its ID
func DeleteSocialMediaEntryByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from token or request context
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	socialMediaIDStr := mux.Vars(r)["socialMediaID"]
	socialMediaID, err := strconv.Atoi(socialMediaIDStr)
	if err != nil {
		http.Error(w, "Invalid social media ID", http.StatusBadRequest)
		return
	}

	// Check if the social media entry exists
	var existingSocialMedia models.SocialMedia
	err = database.GetDB().Where("user_id = ? AND id = ?", userID, socialMediaID).First(&existingSocialMedia).Error
	if err != nil {
		http.Error(w, "Social media entry not found", http.StatusNotFound)
		return
	}

	// Delete the social media entry from the database
	err = database.GetDB().Delete(&existingSocialMedia).Error
	if err != nil {
		http.Error(w, "Failed to delete social media entry", http.StatusInternalServerError)
		return
	}

	// Set appropriate response status
	w.WriteHeader(http.StatusOK)
}
