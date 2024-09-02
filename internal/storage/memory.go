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

type Memory struct {
	mapData sync.Map
	encoder *json.Encoder
}

func NewMemory() *Memory {
	return &Memory{}
}

func (s *Memory) AddOriginalURL(originalURL, userID string) (string, error) {
	var err error = nil
	hash := GenerateHash(originalURL)
	shortData, exist := s.getOrCreateShortData(hash, originalURL, userID)
	if exist {
		err = cerrors.ErrNewShortURL
	}
	return shortData.ShortURL, err
}

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

func (s *Memory) DeleteOriginalURL(ctx context.Context, data service.RequestUserShortenedURL) error {
	err := cerrors.ErrGetShortURLDelete
	s.mapData.Range(func(key, value interface{}) bool {
		safeData := value.(*model.SafeShortData).Load()
		if safeData.ShortURL == data.ShortURL {
			if safeData.UserID == data.UserID {
				shd := safeData
				shd.DeletedFlag = true
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

func (s *Memory) PingContext() error {
	return errors.New("ping MEMORY not yet implemented ")
}

func (s *Memory) Close() error {
	return nil
}

func GenerateHash(url string) string {
	hash := sha256.Sum256([]byte(url))
	return hex.EncodeToString(hash[:])
}

func RandText(size int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") // 52
	b := make([]rune, size)
	for i := range b {
		resI, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterRunes))))
		if err != nil {
			log.Printf("error rand.Int error: %s", err)
			return ""
		}
		b[i] = letterRunes[resI.Int64()]
	}
	return string(b)
}
