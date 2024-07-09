package storage

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/model"
	"log"
	"math/big"
)

type Memory struct {
	mapData map[string]model.ShortData
	encoder *json.Encoder
}

func NewMemory() *Memory {
	return &Memory{
		mapData: make(map[string]model.ShortData),
	}
}

func (s *Memory) GetListShortData() map[string]model.ShortData {
	return s.mapData
}

func (s *Memory) GetURL(shURL string) string {
	for _, value := range s.mapData {
		if shURL == value.ShortURL {
			return value.OriginalURL
		}
	}
	return ""
}

func (s *Memory) AddURL(originalURL string) string {
	hash := GenerateHash(originalURL)
	shortData, _ := s.getOrCreateShortData(hash, originalURL)
	return shortData.ShortURL
}

func (s *Memory) AddBatch(req []shorturl.RequestBatch) []shorturl.ResponseBatch {
	resp := make([]shorturl.ResponseBatch, len(req))
	for i, val := range req {
		sh := s.AddURL(val.OriginalURL)
		resp[i] = shorturl.ResponseBatch{CorrelationID: val.CorrelationID, ShortURL: sh}
	}

	return resp
}

func (s *Memory) getOrCreateShortData(hash, url string) (*model.ShortData, bool) {
	shortData, exist := s.mapData[hash]
	if !exist {
		shortData = model.ShortData{UUID: hash, ShortURL: RandText(8), OriginalURL: url}
		s.mapData[hash] = shortData
	}
	return &shortData, !exist
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
