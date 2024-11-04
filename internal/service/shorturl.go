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

// Storage представляет интерфейс для работы с хранилищем URL-ов.
// Он включает методы для добавления, получения, удаления URL-ов и проверки доступности хранилища.
type Storage interface {
	// AddOriginalURL добавляет новый оригинальный URL, связанный с userID.
	// Возвращает сокращённый URL или ошибку.
	AddOriginalURL(string, string) (string, error)

	// AddOriginalURLBatch добавляет несколько оригинальных URL в пакетном режиме, привязанных к userID.
	// Возвращает срез структур ResponseBatch с сокращёнными URL и корреляционными ID.
	AddOriginalURLBatch([]RequestBatch, string, string) ([]ResponseBatch, error)

	// GetOriginalURL возвращает оригинальный URL по сокращённому, если он существует.
	// Возвращает ошибку, если URL удалён или не найден.
	GetOriginalURL(context.Context, string) (model.ShortData, error)

	// GetUserShortenedURLs возвращает список всех сокращённых URL, созданных пользователем с указанным userID.
	GetUserShortenedURLs(context.Context, string, string) ([]ResponseUserShortenedURL, error)

	// GetURLs возвращает количество сокращённых URL в сервисе
	GetURLs() int

	// DeleteOriginalURL удаляет сокращённый URL по данным пользователя.
	DeleteOriginalURL(context.Context, RequestUserShortenedURL) error

	// PingContext проверяет доступность хранилища.
	PingContext() error

	// Close закрывает соединение с хранилищем.
	Close() error
}

// ShortURL хранит информацию о пользователях системы и их текущем состоянии авторизации.
type ShortURL struct {
	RedirectURL string  // базовый URL на который производится редирект
	storage     Storage // хранилище URL-ов
}

// RequestBatch представляет структуру запроса для пакетного добавления URL.
type RequestBatch struct {
	CorrelationID string `json:"correlation_id"` // уникальный ID для идентификации запроса
	OriginalURL   string `json:"original_url"`   // оригинальный URL для сокращения
}

// ResponseBatch представляет структуру ответа для пакетного добавления URL.
type ResponseBatch struct {
	CorrelationID string `json:"correlation_id"` // уникальный ID запроса
	ShortURL      string `json:"short_url"`      // сокращённый URL
}

// ResponseUserShortenedURL содержит информацию о сокращённом URL и оригинальном URL пользователя.
type ResponseUserShortenedURL struct {
	OriginalURL string `json:"original_url"` // оригинальный URL
	ShortURL    string `json:"short_url"`    // сокращённый URL
}

// RequestUserShortenedURL содержит информацию о пользователе и сокращённом URL для удаления.
type RequestUserShortenedURL struct {
	UserID   string `json:"user_id"`   // ID пользователя
	ShortURL string `json:"short_url"` // сокращённый URL для удаления
}

// NewShortURL создаёт новый экземпляр ShortURL с базовым URL для редиректов и хранилищем для управления URL.
func NewShortURL(redirectURL string, strg Storage) *ShortURL {
	s := new(ShortURL)
	s.RedirectURL = redirectURL
	s.storage = strg
	return s
}

// SetRedirectURL изменяет базовый URL для редиректа.
func (s *ShortURL) SetRedirectURL(url string) {
	s.RedirectURL = url
}

// AddOriginalURL добавляет оригинальный URL и возвращает полный сокращённый URL, используя базовый URL редиректа.
func (s *ShortURL) AddOriginalURL(url, userID string) (string, error) {
	res, err := s.storage.AddOriginalURL(url, userID)
	return fmt.Sprintf("%s/%s", s.RedirectURL, res), err
}

// AddOriginalURLBatch добавляет несколько оригинальных URL-адресов в пакетном режиме для указанного пользователя.
func (s *ShortURL) AddOriginalURLBatch(batch []RequestBatch, userID string) ([]ResponseBatch, error) {
	return s.storage.AddOriginalURLBatch(batch, s.RedirectURL, userID)
}

// GetOriginalURL возвращает оригинальный URL по сокращённому. Если URL был удалён, возвращает ошибку.
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

// GetUserShortenedURLs возвращает список всех сокращённых URL пользователя в формате JSON.
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

// GetURLs возвращает количество сокращённых URL в сервисе
func (s *ShortURL) GetURLs() int {
	return s.storage.GetURLs()
}

// DeleteOriginalURLs удаляет список сокращённых URL-адресов для пользователя в асинхронном режиме с использованием горутин.
func (s *ShortURL) DeleteOriginalURLs(ctx context.Context, userID string, shortURLs []string) {
	inCH := s.generatorDeleteShortURL(userID, shortURLs)
	channels := s.funOut(inCH)
	s.funIn(channels...)
}

// PingContext проверяет доступность хранилища.
func (s *ShortURL) PingContext() error {
	return s.storage.PingContext()
}

// Close закрывает соединение с хранилищем.
func (s *ShortURL) Close() error {
	err := s.storage.Close()
	return err
}

// generatorDeleteShortURL создаёт канал, через который передаются URL для удаления.
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

// deleteOriginalURL удаляет сокращённый URL из хранилища, используя данные, полученные из канала.
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

// funIn собирает результаты удаления URL-адресов в один финальный канал.
// Использует WaitGroup для синхронизации завершения горутин.
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

// funOut создаёт и возвращает несколько каналов для асинхронного удаления URL-адресов.
func (s *ShortURL) funOut(inCH chan RequestUserShortenedURL) []chan error {

	numWorkers := 20
	channels := make([]chan error, numWorkers)

	for i := 0; i < numWorkers; i++ {
		addResultCh := s.deleteOriginalURL(inCH)
		channels[i] = addResultCh
	}

	return channels
}
