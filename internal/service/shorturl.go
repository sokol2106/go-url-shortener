//go:generate mockgen -destination=../../../mocks/mock_shorturl.go -package=shorturl . StorageURL
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
)

type Storage interface {
	AddOriginalURL(string, string) (string, error)
	AddOriginalURLBatch([]RequestBatch, string, string) ([]ResponseBatch, error)
	GetOriginalURL(context.Context, string) string
	GetUserShortenedURLs(context.Context, string, string) ([]ResponseUserShortenedURL, error)
	PingContext() error
	Close() error
}

type ShortURL struct {
	RedirectURL string
	storage     Storage
}

type RequestBatch struct {
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}

type ResponseUserShortenedURL struct {
	OriginalURL string `json:"original_url"`
	ShortURL    string `json:"short_url"`
}

func NewShortURL(redirectURL string, strg Storage) *ShortURL {
	s := new(ShortURL)
	s.RedirectURL = redirectURL
	s.storage = strg
	return s
}

func (s *ShortURL) AddOriginalURL(url, userID string) (string, error) {
	res, err := s.storage.AddOriginalURL(url, userID)
	return fmt.Sprintf("%s/%s", s.RedirectURL, res), err
}

func (s *ShortURL) AddOriginalURLBatch(batch []RequestBatch, userID string) ([]ResponseBatch, error) {
	return s.storage.AddOriginalURLBatch(batch, s.RedirectURL, userID)
}

func (s *ShortURL) GetOriginalURL(ctx context.Context, url string) string {
	return s.storage.GetOriginalURL(ctx, url)
}

func (s *ShortURL) GetUserShortenedURLs(ctx context.Context, userID string) ([]byte, error) {
	var res bytes.Buffer
	ctx2, cancel := context.WithCancel(ctx)
	defer cancel()
	responseBatch, err := s.storage.GetUserShortenedURLs(ctx2, userID, s.RedirectURL)
	if err != nil {
		return nil, err
	}

	err = json.NewEncoder(&res).Encode(responseBatch)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(&res)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (s *ShortURL) PingContext() error {
	return s.storage.PingContext()
}

func (s *ShortURL) Close() error {
	err := s.storage.Close()
	return err
}
