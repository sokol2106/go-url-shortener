package storage

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sokol2106/go-url-shortener/internal/cerrors"
	"github.com/sokol2106/go-url-shortener/internal/model"
	"github.com/sokol2106/go-url-shortener/internal/service"
	"log"
	"math/big"
	"sync"
)

// Memory представляет структуру для работы со структурой Map, реализующую интерфейс Storage.
// mapData - это хранилище оригинальных и сокращённых URL,
// encoder - для записи, data хранит считанные из файла данные.
type Memory struct {
	mapData sync.Map      // потокобезопасная структура для хранения URL
	encoder *json.Encoder // энкодер JSON для записи данных в файл
}

// NewMemory создаёт объект Memory и возвращает ссылку на него
func NewMemory() *Memory {
	return &Memory{}
}

// AddOriginalURL добавляет оригинальный URL и возвращает сгенерированный сокращённый URL.
func (s *Memory) AddOriginalURL(originalURL, userID string) (string, error) {
	var err error = nil
	hash := GenerateHash(originalURL)
	shortData, exist := s.getOrCreateShortData(hash, originalURL, userID)
	if exist {
		err = cerrors.ErrNewShortURL
	}
	return shortData.ShortURL, err
}

// AddOriginalURLBatch добавляет список URL-ов в пакетном режиме и возвращает соответствующие сокращённые URL-ы.
func (s *Memory) AddOriginalURLBatch(req []service.RequestBatch, redirectURL string, userID string) ([]service.ResponseBatch, error) {
	var err error = nil
	resp := make([]service.ResponseBatch, len(req))
	for i, val := range req {
		sh, addErr := s.AddOriginalURL(val.OriginalURL, userID)
		if addErr != nil {
			err = addErr
		}
		resp[i] = service.ResponseBatch{CorrelationID: val.CorrelationID, ShortURL: fmt.Sprintf("%s/%s", redirectURL, sh)}
	}

	return resp, err
}

// GetOriginalURL ищет оригинальный URL по сокращённому. Если URL не найден, возвращает ошибку.
func (s *Memory) GetOriginalURL(ctx context.Context, shURL string) (model.ShortData, error) {
	var (
		result model.ShortData
		err    error = nil
	)
	err = cerrors.ErrGetShortURLNotFind
	s.mapData.Range(func(key, value interface{}) bool {
		safeData := value.(*model.SafeShortData).Load()
		if shURL == safeData.ShortURL {
			result = safeData
			err = nil
			return false
		}
		return true
	})

	return result, err
}

// GetUserShortenedURLs возвращает список сокращённых URL-ов пользователя, найденных в файле.
func (s *Memory) GetUserShortenedURLs(ctx context.Context, userID, redirectURL string) ([]service.ResponseUserShortenedURL, error) {
	result := make([]service.ResponseUserShortenedURL, 0)
	s.mapData.Range(func(key, value interface{}) bool {
		safeData := value.(*model.SafeShortData).Load()
		if userID == safeData.UserID {
			result = append(result, service.ResponseUserShortenedURL{OriginalURL: safeData.OriginalURL, ShortURL: fmt.Sprintf("%s/%s", redirectURL, safeData.ShortURL)})
		}
		return true
	})

	return result, nil
}

// DeleteOriginalURL не реализован. Возвращает nil.
func (s *Memory) DeleteOriginalURL(ctx context.Context, data service.RequestUserShortenedURL) error {
	err := cerrors.ErrGetShortURLDelete
	s.mapData.Range(func(key, value interface{}) bool {
		safeData := value.(*model.SafeShortData).Load()
		if safeData.ShortURL == data.ShortURL {
			if safeData.UserID == data.UserID {
				//shd := safeData
				safeData.DeletedFlag = true
				value.(*model.SafeShortData).Store(safeData)
				//s.mapData.Store(key, shd)
				err = nil
			}
			return false
		}
		return true
	})
	return err
}

// getOrCreateShortData ищет данные о сокращённом URL или создаёт их, если они не найдены.
// Возвращает найденные или созданные данные, а также флаг, указывающий, были ли данные новыми.
func (s *Memory) getOrCreateShortData(hash, url, userID string) (*model.ShortData, bool) {
	var shortData model.ShortData
	value, exist := s.mapData.Load(hash)
	if exist {
		shortData = value.(*model.SafeShortData).Load()
	} else {
		shortData = model.ShortData{UUID: hash, ShortURL: RandText(8), OriginalURL: url, UserID: userID, DeletedFlag: false}
		s.mapData.Store(hash, model.NewSafeShortData(shortData))
	}
	return &shortData, exist
}

// PingContext не реализован. Возвращает ошибку.
func (s *Memory) PingContext() error {
	return errors.New("ping MEMORY not yet implemented ")
}

// Close не реализован. Возвращает nil.
func (s *Memory) Close() error {
	return nil
}

// GenerateHash рассчитывает hash URL на основе sha256
func GenerateHash(url string) string {
	hash := sha256.Sum256([]byte(url))
	return hex.EncodeToString(hash[:])
}

// RandText возвращает рандомизированный текст, размером size.
// Используется библиотека "crypto/rand".
func RandText(size int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") // 52
	var b []rune
	//	b := make([]rune, size)
	for i := 0; i < size; i++ {
		resI, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			log.Printf("error rand.Int error: %s", err)
			return ""
		}
		b = append(b, letterRunes[resI.Int64()])
		//	b[i] = letterRunes[resI.Int64()]
	}
	return string(b)
}
