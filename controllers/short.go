package controllers

import (
	models "github.com/yatiac/go-shortener/models"
	storage "github.com/yatiac/go-shortener/storage"
)

type ShortController struct {
	store *storage.DBStore
}

func NewShortController(store *storage.DBStore) *ShortController {
	return &ShortController{store: store}
}
func (sc *ShortController) CreateShortURL(longURL, slug string) (*models.URL, error) {
	url := &models.URL{
		LongURL: longURL,
		Slug:    slug,
	}
	if err := sc.store.CreateURL(url); err != nil {
		return nil, err
	}
	return url, nil
}
