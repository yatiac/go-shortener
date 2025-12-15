package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	services "github.com/yatiac/go-shortener/services"
)

type ShortController struct {
	service *services.ShortenerService
}

func NewShortController(service *services.ShortenerService) *ShortController {
	return &ShortController{service: service}
}
func (sc *ShortController) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	var request struct {
		LongURL string `json:"long_url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if request.LongURL == "" {
		http.Error(w, "Long URL is required", http.StatusBadRequest)
		return
	}
	url, err := sc.service.CreateShortURL(request.LongURL)

	if err != nil {
		http.Error(w, "Failed to create short URL: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		ShortURL string `json:"short_url"`
		Slug     string `json:"slug"`
		LongURL  string `json:"long_url"`
	}{
		ShortURL: os.Getenv("BASE_URL") + "/" + url.Slug,
		Slug:     url.Slug,
		LongURL:  url.LongURL,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func (sc *ShortController) RedirectToLongURL(w http.ResponseWriter, r *http.Request) {
	slug := r.URL.Path[1:] // Remove leading slash

	url, err := sc.service.GetLongURL(slug)
	if err != nil {
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, url.LongURL, http.StatusMovedPermanently)
}
