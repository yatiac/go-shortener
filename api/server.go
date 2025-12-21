package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/yatiac/go-shortener/controllers"
	"github.com/yatiac/go-shortener/models"
	"github.com/yatiac/go-shortener/services"
	"github.com/yatiac/go-shortener/storage"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Handler(w http.ResponseWriter, r *http.Request) {
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
	router := mux.NewRouter()

	// API routes
	router.HandleFunc("/api/shorten", shortController.CreateShortURL).Methods(http.MethodPost, http.MethodOptions)

	router.HandleFunc("/kaithheathcheck", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods(http.MethodGet)

	router.HandleFunc("/list_files", func(w http.ResponseWriter, r *http.Request) {
		files, err := os.ReadDir("./web/dist")
		if err != nil {
			http.Error(w, "Failed to read directory", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("["))
		for i, file := range files {
			w.Write([]byte(fmt.Sprintf("\"%s\"", file.Name())))
			if i < len(files)-1 {
				w.Write([]byte(","))
			}
		}
		w.Write([]byte("]"))
	}).Methods(http.MethodGet)

	// Return long URL
	router.HandleFunc("/api/short_url/{slug}", shortController.GetLongURL).Methods(http.MethodGet)

	// Serve static files from web/dist
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./web/dist/assets"))))

	// Catch-all handler: serve index.html for SPA routing (must be last)
	router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		http.ServeFile(w, req, "./web/dist/index.html")
	})

	// CORS
	router.Use(controllers.DisableCORS)
	router.Use(mux.CORSMethodMiddleware(router))

	// Start server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	router.ServeHTTP(w, r)
}
