package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/faazabilamri7/mygram/database"
	"github.com/faazabilamri7/mygram/models"
	"github.com/gorilla/mux"
)

// CreateComment handles the creation of a new comment
func CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment models.Comment
	err := json.NewDecoder(r.Body).Decode(&comment)
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

	// Set user ID for the comment
	comment.UserID = userID

	// Save comment to the database
	err = database.GetDB().Create(&comment).Error
	if err != nil {
		http.Error(w, "Failed to create comment", http.StatusInternalServerError)
		return
	}

	// Set appropriate response status and return the created comment
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(comment)
}

// GetAllComments handles fetching all comments from all users
func GetAllComments(w http.ResponseWriter, r *http.Request) {
	var comments []models.Comment
	err := database.GetDB().Preload("User").Preload("Photo").Find(&comments).Error
	if err != nil {
		http.Error(w, "Failed to fetch comments", http.StatusInternalServerError)
		return
	}

	// Set appropriate response status and return the fetched comments
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comments)
}

// GetCommentByID handles fetching a comment by its ID
func GetCommentByID(w http.ResponseWriter, r *http.Request) {
	commentIDStr := mux.Vars(r)["commentId"]
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var comment models.Comment
	err = database.GetDB().Preload("User").Preload("Photo").First(&comment, commentID).Error
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	// Set appropriate response status and return the fetched comment
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(comment)
}

// UpdateCommentByID handles updating a comment by its ID
func UpdateCommentByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from token or request context
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	commentIDStr := mux.Vars(r)["commentId"]
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	var updatedComment models.Comment
	err = json.NewDecoder(r.Body).Decode(&updatedComment)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Check if the comment exists
	var existingComment models.Comment
	err = database.GetDB().First(&existingComment, commentID).Error
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	// Check if the user is authorized to update the comment
	if existingComment.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Update comment message
	existingComment.Message = updatedComment.Message

	// Save updated comment to the database
	err = database.GetDB().Save(&existingComment).Error
	if err != nil {
		http.Error(w, "Failed to update comment", http.StatusInternalServerError)
		return
	}

	// Set appropriate response status and return the updated comment
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingComment)
}

// DeleteCommentByID handles deleting a comment by its ID
func DeleteCommentByID(w http.ResponseWriter, r *http.Request) {
	// Retrieve user ID from token or request context
	userID, err := getUserIDFromToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	commentIDStr := mux.Vars(r)["commentId"]
	commentID, err := strconv.Atoi(commentIDStr)
	if err != nil {
		http.Error(w, "Invalid comment ID", http.StatusBadRequest)
		return
	}

	// Check if the comment exists
	var existingComment models.Comment
	err = database.GetDB().First(&existingComment, commentID).Error
	if err != nil {
		http.Error(w, "Comment not found", http.StatusNotFound)
		return
	}

	// Check if the user is authorized to delete the comment
	if existingComment.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Delete the comment from the database
	err = database.GetDB().Delete(&existingComment).Error
	if err != nil {
		http.Error(w, "Failed to delete comment", http.StatusInternalServerError)
		return
	}

	// Set appropriate response status
	w.WriteHeader(http.StatusOK)
}
