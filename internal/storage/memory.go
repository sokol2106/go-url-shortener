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
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/model"
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

func (s *Memory) GetURL(ctx context.Context, shURL string) string {
	original := ""
	s.mapData.Range(func(key, value interface{}) bool {
		mdl := value.(model.ShortData)
		if shURL == mdl.ShortURL {
			original = mdl.OriginalURL
			return false
		}
		return true
	})
	return original
}

func (s *Memory) AddURL(originalURL string) (string, error) {
	var err error = nil
	hash := GenerateHash(originalURL)
	shortData, exist := s.getOrCreateShortData(hash, originalURL)
	if exist {
		err = cerrors.ErrNewShortURL
	}
	return shortData.ShortURL, err
}

func (s *Memory) AddBatch(req []shorturl.RequestBatch, redirectURL string) ([]shorturl.ResponseBatch, error) {
	var err error = nil
	resp := make([]shorturl.ResponseBatch, len(req))
	for i, val := range req {
		sh, addErr := s.AddURL(val.OriginalURL)
		if addErr != nil {
			err = addErr
		}
		resp[i] = shorturl.ResponseBatch{CorrelationID: val.CorrelationID, ShortURL: fmt.Sprintf("%s/%s", redirectURL, sh)}
	}

	return resp, err
}

func (s *Memory) getOrCreateShortData(hash, url string) (*model.ShortData, bool) {
	var shortData model.ShortData
	value, exist := s.mapData.Load(hash)
	if exist {
		shortData = value.(model.ShortData)
	} else {
		shortData = model.ShortData{UUID: hash, ShortURL: RandText(8), OriginalURL: url}
		s.mapData.Store(hash, shortData)
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
