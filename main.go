package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ArtemZ007/em-test/config"
	"github.com/ArtemZ007/em-test/internal/handlers"
	"github.com/ArtemZ007/em-test/internal/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Migrate the database
	db.AutoMigrate(&models.Song{})

	// Initialize the handler
	songHandler := handlers.SongHandler{DB: db}

	// Create router and define routes
	r := mux.NewRouter()

	r.HandleFunc("/songs", songHandler.GetSongs).Methods("GET")
	r.HandleFunc("/songs", songHandler.AddSong).Methods("POST")
	r.HandleFunc("/songs/{id}", songHandler.UpdateSong).Methods("PUT")
	r.HandleFunc("/songs/{id}", songHandler.DeleteSong).Methods("DELETE")

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
