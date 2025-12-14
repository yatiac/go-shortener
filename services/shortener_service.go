package services

import (
	"math/rand"
	"time"

	_ "github.com/joho/godotenv/autoload"
	models "github.com/yatiac/go-shortener/models"
	storage "github.com/yatiac/go-shortener/storage"
)

type ShortenerService struct {
	store *storage.DBStore
}

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewShortenerService(store *storage.DBStore) *ShortenerService {
	return &ShortenerService{store: store}
}

func (s *ShortenerService) CreateShortURL(longURL string) (*models.URL, error) {
	// Create a unique slug
	slug := generateSlug()

	url := &models.URL{
		LongURL: longURL,
		Slug:    slug,
	}

	if err := s.store.CreateURL(url); err != nil {
		return nil, err
	}
	return url, nil
}

func generateSlug() string {
	b := make([]byte, 9)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (s *ShortenerService) GetLongURL(slug string) (*models.URL, error) {
	return s.store.GetURLBySlug(slug)
}
