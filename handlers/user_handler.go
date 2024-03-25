package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/faazabilamri7/mygram/database"
	"github.com/faazabilamri7/mygram/models"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser handles the registration of a new user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Hash the password before saving it to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)

	// Set created and updated timestamps
	currentTime := time.Now()
	user.CreatedAt = currentTime
	user.UpdatedAt = currentTime

	// Save user to database using GORM
	err = database.GetDB().Create(&user).Error
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		return
	}

	// Set appropriate response status and return the registered user
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// LoginUser handles user login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Decode request body
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if user with provided email exists
	var user models.User
	if err := database.GetDB().Where("email = ?", credentials.Email).First(&user).Error; err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Check if the password is correct
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password)); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := generateToken(int64(user.ID), user.Email)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Set appropriate response status and return the token
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"token": token,
	})
}

// UpdateUser handles updating user information
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	var updateUserReq models.User
	err := json.NewDecoder(r.Body).Decode(&updateUserReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Retrieve user ID from the JWT token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Retrieve existing user from the database
	var user models.User
	err = database.GetDB().First(&user, userID).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Update user fields based on the request
	user.Username = updateUserReq.Username
	user.Email = updateUserReq.Email
	user.Age = updateUserReq.Age
	user.ImageURL = updateUserReq.ImageURL

	// Save updated user to the database
	err = database.GetDB().Save(&user).Error
	if err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	// Return updated user as response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// DeleteUser handles deleting a user account
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from the JWT token
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Retrieve existing user from the database
	var user models.User
	err = database.GetDB().First(&user, userID).Error
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Delete user from the database
	err = database.GetDB().Delete(&user).Error
	if err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	// Return success message
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "User deleted successfully"})
}
