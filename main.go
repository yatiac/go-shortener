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
	r.HandleFunc("/api/shorten", shortController.CreateShortURL).Methods(http.MethodPost, http.MethodOptions)

	r.HandleFunc("/kaithheathcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	// Return long URL
	r.HandleFunc("/{slug}", shortController.GetLongURL).Methods(http.MethodGet)

	// CORS
	r.Use(disableCORS)
	r.Use(mux.CORSMethodMiddleware(r))

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func disableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
