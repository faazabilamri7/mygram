package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/faazabilamri7/mygram/database"
	"github.com/faazabilamri7/mygram/handlers"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize database connection
	database.InitDB()

	// Migrate database schema
	database.AutoMigrate()

	// Check database connection
	err := database.GetDB().DB().Ping()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	r := mux.NewRouter()

	//Welcome
	r.HandleFunc("/", welcomeMessage).Methods("GET")

	// Register API endpoints
	r.HandleFunc("/users/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/users/login", handlers.LoginUser).Methods("POST")
	r.HandleFunc("/users", handlers.UpdateUser).Methods("PUT")
	r.HandleFunc("/users", handlers.DeleteUser).Methods("DELETE")

	r.HandleFunc("/photos", handlers.CreatePhoto).Methods("POST")
	r.HandleFunc("/photos", handlers.GetAllPhotos).Methods("GET")
	r.HandleFunc("/photos/{photoID}", handlers.GetPhotoByID).Methods("GET")
	r.HandleFunc("/photos/{photoID}", handlers.UpdatePhotoByID).Methods("PUT")
	r.HandleFunc("/photos/{photoID}", handlers.DeletePhotoByID).Methods("DELETE")

	r.HandleFunc("/comments", handlers.CreateComment).Methods("POST")
	r.HandleFunc("/comments", handlers.GetAllComments).Methods("GET")
	r.HandleFunc("/comments/{commentID}", handlers.GetCommentByID).Methods("GET")
	r.HandleFunc("/comments/{commentID}", handlers.UpdateCommentByID).Methods("PUT")
	r.HandleFunc("/comments/{commentID}", handlers.DeleteCommentByID).Methods("DELETE")

	r.HandleFunc("/socialmedias", handlers.CreateSocialMediaEntry).Methods("POST")
	r.HandleFunc("/socialmedias", handlers.GetAllSocialMediaEntries).Methods("GET")
	r.HandleFunc("/socialmedias/{socialMediaID}", handlers.GetSocialMediaEntryByID).Methods("GET")
	r.HandleFunc("/socialmedias/{socialMediaID}", handlers.UpdateSocialMediaEntryByID).Methods("PUT")
	r.HandleFunc("/socialmedias/{socialMediaID}", handlers.DeleteSocialMediaEntryByID).Methods("DELETE")

	// Start server
	fmt.Println("Server started at localhost:443")
	log.Fatal(http.ListenAndServeTLS(":443", "cert.pem", "key.pem", r))
}

func welcomeMessage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Selamat datang, myGram app Faaza bil amri, Golang 004, Newbie")
}
