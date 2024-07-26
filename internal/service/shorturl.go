//go:generate mockgen -destination=../../../mocks/mock_shorturl.go -package=shorturl . StorageURL
package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/sokol2106/go-url-shortener/internal/cerrors"
	"github.com/sokol2106/go-url-shortener/internal/model"
	"io"
	"sync"
)

type Storage interface {
	AddOriginalURL(string, string) (string, error)
	AddOriginalURLBatch([]RequestBatch, string, string) ([]ResponseBatch, error)
	GetOriginalURL(context.Context, string) (model.ShortData, error)
	GetUserShortenedURLs(context.Context, string, string) ([]ResponseUserShortenedURL, error)
	DeleteOriginalURL(context.Context, RequestUserShortenedURL) error
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

type RequestUserShortenedURL struct {
	UserID   string `json:"user_id"`
	ShortURL string `json:"short_url"`
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

func (s *ShortURL) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	mdl, err := s.storage.GetOriginalURL(ctx, shortURL)
	if err != nil {
		return "", err
	}

	if mdl.DeletedFlag {
		return "", cerrors.ErrGetShortURLDelete
	}

	return mdl.OriginalURL, nil
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

func (s *ShortURL) DeleteOriginalURLs(ctx context.Context, userID string, shortURLs []string) {
	inCH := s.generatorDeleteShortURL(userID, shortURLs)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	s.deleteOriginalURL(inCH)
	//resultCh := s.funIn(ch1, ch2, ch3, ch4, ch5, ch6, ch7, ch8, ch9, ch10)

	//for res := range resultCh {
	//	fmt.Println(res)
	//}
}

func (s *ShortURL) PingContext() error {
	return s.storage.PingContext()
}

func (s *ShortURL) Close() error {
	err := s.storage.Close()
	return err
}

func (s *ShortURL) generatorDeleteShortURL(userID string, shortURLs []string) chan RequestUserShortenedURL {
	inputCh := make(chan RequestUserShortenedURL)

	go func() {
		defer close(inputCh)
		for _, value := range shortURLs {
			data := RequestUserShortenedURL{UserID: userID, ShortURL: value}
			inputCh <- data
		}
	}()

	return inputCh
}

func (s *ShortURL) deleteOriginalURL(inputCh chan RequestUserShortenedURL) chan error {
	resultCh := make(chan error)

	go func() {
		defer close(resultCh)
		for data := range inputCh {
			resultCh <- s.storage.DeleteOriginalURL(context.Background(), data)
		}
	}()

	return resultCh
}

func (s *ShortURL) funIn(chs ...chan error) chan error {
	finalCh := make(chan error)
	var wg sync.WaitGroup

	for _, ch := range chs {
		chClosure := ch
		wg.Add(1)

		go func() {
			defer wg.Done()
			for data := range chClosure {
				finalCh <- data
			}
		}()
	}

	go func() {
		wg.Wait()
		close(finalCh)
	}()

	return finalCh
}
