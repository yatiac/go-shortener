package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv"
	controllers "github.com/yatiac/go-shortener/controllers"
	models "github.com/yatiac/go-shortener/models"
	services "github.com/yatiac/go-shortener/services"
	storage "github.com/yatiac/go-shortener/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	if err := db.AutoMigrate(&models.URL{}); err != nil {
		log.Fatalf("Failed to auto-migrate models: %v", err)
	}

	dbStore := storage.NewDBStore(db)
	shortService := services.NewShortenerService(dbStore)
	shortController := controllers.NewShortController(shortService)

	// Set up Gorilla Mux router
	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/shorten", shortController.CreateShortURL).Methods("POST")

	// Redirect route (should be last)
	r.HandleFunc("/{slug}", shortController.RedirectToLongURL).Methods("GET")

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
