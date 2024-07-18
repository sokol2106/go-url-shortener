//go:generate mockgen -destination=../../../mocks/mock_shorturl.go -package=shorturl . StorageURL
package service

import (
	"context"
	"fmt"
)

type StorageURL interface {
	AddURL(string) (string, error)
	AddBatch([]RequestBatch, string) ([]ResponseBatch, error)
	GetURL(context.Context, string) string
	PingContext() error
	Close() error
}

type ShortURL struct {
	RedirectURL string
	storageURL  StorageURL
}

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
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

func (s *ShortURL) AddURL(url string) (string, error) {
	res, err := s.storageURL.AddURL(url)
	return fmt.Sprintf("%s/%s", s.RedirectURL, res), err
}

func (s *ShortURL) AddBatch(batch []RequestBatch) ([]ResponseBatch, error) {
	return s.storageURL.AddBatch(batch, s.RedirectURL)
}

func (s *ShortURL) GetURL(ctx context.Context, url string) string {
	return s.storageURL.GetURL(ctx, url)
}

func (s *ShortURL) PingContext() error {
	return s.storageURL.PingContext()
}

func (s *ShortURL) Close() error {
	err := s.storageURL.Close()
	return err
}
