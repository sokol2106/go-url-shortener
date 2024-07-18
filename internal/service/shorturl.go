//go:generate mockgen -destination=../../../mocks/mock_shorturl.go -package=shorturl . StorageURL
package service

import (
	"context"
	"fmt"
)

type StorageURL interface {
	AddOriginalURL(string, string) (string, error)
	AddOriginalURLBatch([]RequestBatch, string, string) ([]ResponseBatch, error)
	GetOriginalURL(context.Context, string) string
	PingContext() error
	Close() error
}

type ShortURL struct {
	RedirectURL string
	storageURL  StorageURL
}

type RequestBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

func NewShortURL(redirectURL string, strg StorageURL) *ShortURL {
	s := new(ShortURL)
	s.RedirectURL = redirectURL
	s.storageURL = strg
	return s
}

func (s *ShortURL) AddURL(url, userID string) (string, error) {
	res, err := s.storageURL.AddOriginalURL(url, userID)
	return fmt.Sprintf("%s/%s", s.RedirectURL, res), err
}

func (s *ShortURL) AddBatch(batch []RequestBatch, userID string) ([]ResponseBatch, error) {
	return s.storageURL.AddOriginalURLBatch(batch, s.RedirectURL, userID)
}

func (s *ShortURL) GetURL(ctx context.Context, url string) string {
	return s.storageURL.GetOriginalURL(ctx, url)
}

func (s *ShortURL) PingContext() error {
	return s.storageURL.PingContext()
}

func (s *ShortURL) Close() error {
	err := s.storageURL.Close()
	return err
}
