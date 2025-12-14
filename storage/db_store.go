package storage

import (
	"github.com/yatiac/go-shortener/models"
	"gorm.io/gorm"
)

type DBStore struct {
	DB *gorm.DB
}

func NewDBStore(db *gorm.DB) *DBStore {
	return &DBStore{DB: db}
}

func (store *DBStore) CreateURL(url *models.URL) error {
	return store.DB.Create(url).Error
}

func (store *DBStore) GetURLBySlug(slug string) (*models.URL, error) {
	var url models.URL
	if err := store.DB.Where("slug = ?", slug).First(&url).Error; err != nil {
		return nil, err
	}
	return &url, nil
}
