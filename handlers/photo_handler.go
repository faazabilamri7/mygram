package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/faazabilamri7/mygram/database"
	"github.com/faazabilamri7/mygram/models"
	"github.com/gorilla/mux"
)

// CreatePhoto handles the creation of a new photo
func CreatePhoto(w http.ResponseWriter, r *http.Request) {
	var photo models.Photo
	err := json.NewDecoder(r.Body).Decode(&photo)
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

	// Set user ID for the photo
	photo.UserID = userID

	// Save photo to the database
	err = database.GetDB().Create(&photo).Error
	if err != nil {
		http.Error(w, "Failed to create photo", http.StatusInternalServerError)
		return
	}

	// Set appropriate response status and return the created photo
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(photo)
}

// GetAllPhotos handles fetching all photos from all users
func GetAllPhotos(w http.ResponseWriter, r *http.Request) {
	var photos []models.Photo
	err := database.GetDB().Preload("User").Find(&photos).Error
	if err != nil {
		http.Error(w, "Failed to fetch photos", http.StatusInternalServerError)
		return
	}

	// Set appropriate response status and return the fetched photos
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(photos)
}

// GetPhotoByID handles fetching a photo by its ID
func GetPhotoByID(w http.ResponseWriter, r *http.Request) {
	photoIDStr := mux.Vars(r)["photoId"]
	photoID, err := strconv.Atoi(photoIDStr)
	if err != nil {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		return
	}

	var photo models.Photo
	err = database.GetDB().Preload("User").First(&photo, photoID).Error
	if err != nil {
		http.Error(w, "Photo not found", http.StatusNotFound)
		return
	}

	// Set appropriate response status and return the fetched photo
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(photo)
}

// UpdatePhotoByID handles updating a photo by its ID
func UpdatePhotoByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from token or request context
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	photoIDStr := mux.Vars(r)["photoId"]
	photoID, err := strconv.Atoi(photoIDStr)
	if err != nil {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		return
	}

	var updatedPhoto models.Photo
	err = json.NewDecoder(r.Body).Decode(&updatedPhoto)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the photo exists
	var existingPhoto models.Photo
	err = database.GetDB().First(&existingPhoto, photoID).Error
	if err != nil {
		http.Error(w, "Photo not found", http.StatusNotFound)
		return
	}

	// Check if the user is authorized to update the photo
	if existingPhoto.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Update photo fields
	existingPhoto.Title = updatedPhoto.Title
	existingPhoto.Caption = updatedPhoto.Caption
	existingPhoto.URL = updatedPhoto.URL

	// Save updated photo to the database
	err = database.GetDB().Save(&existingPhoto).Error
	if err != nil {
		http.Error(w, "Failed to update photo", http.StatusInternalServerError)
		return
	}

	// Set appropriate response status and return the updated photo
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingPhoto)
}

// DeletePhotoByID handles deleting a photo by its ID
func DeletePhotoByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from token or request context
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	photoIDStr := mux.Vars(r)["photoId"]
	photoID, err := strconv.Atoi(photoIDStr)
	if err != nil {
		http.Error(w, "Invalid photo ID", http.StatusBadRequest)
		return
	}

	// Check if the photo exists
	var existingPhoto models.Photo
	err = database.GetDB().First(&existingPhoto, photoID).Error
	if err != nil {
		http.Error(w, "Photo not found", http.StatusNotFound)
		return
	}

	// Check if the user is authorized to delete the photo
	if existingPhoto.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Delete the photo from the database
	err = database.GetDB().Delete(&existingPhoto).Error
	if err != nil {
		http.Error(w, "Failed to delete photo", http.StatusInternalServerError)
		return
	}

	// Set appropriate response status
	w.WriteHeader(http.StatusOK)
}
