package storage

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/sokol2106/go-url-shortener/internal/model"
	"log"
	"math/big"
	"os"
)

// ShortDataList

func (s *ShortDataList) Init(filename string) {
	s.mapData = make(map[string]model.ShortData)
	s.isWriteEnable = false
	if filename != "" {
		newFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("error load file filename: %s , error: %s", filename, err)
			return
		}

		s.isWriteEnable = true
		s.encoder = json.NewEncoder(newFile)
		s.file = newFile
		s.loadDataFile()
	}
}

func (s *ShortDataList) Close() error {
	return s.file.Close()
}

func (s *ShortDataList) loadDataFile() {
	scanner := bufio.NewScanner(s.file)
	for scanner.Scan() {
		data := scanner.Bytes()
		sd := model.ShortData{}
		err := json.Unmarshal(data, &sd)
		if err != nil {
			fmt.Printf("error Unmarshal file error: %s", err)
			return
		}
		s.mapData[sd.UUID] = sd
		sd = model.ShortData{}
	}
}

func (s *ShortDataList) GetListShortData() map[string]model.ShortData {
	return s.mapData
}

func (s *ShortDataList) GetURL(shURL string) string {
	for _, value := range s.mapData {
		if shURL == value.ShortURL {
			return value.OriginalURL
		}
	}
	return ""
}

func (s *ShortDataList) AddURL(originalURL string) string {
	hash := GenerateHash(originalURL)
	shortData, isNewShortData := s.getOrCreateShortData(hash, originalURL)
	if isNewShortData && s.isWriteEnable {
		s.writeToFile(shortData)
	}

	return shortData.ShortURL
}

func (s *ShortDataList) getOrCreateShortData(hash, url string) (*model.ShortData, bool) {
	shortData, exist := s.mapData[hash]
	if !exist {
		shortData = model.ShortData{UUID: hash, ShortURL: RandText(8), OriginalURL: url}
		s.mapData[hash] = shortData
	}
	return &shortData, !exist
}

func (s *ShortDataList) writeToFile(data *model.ShortData) {
	if err := s.encoder.Encode(&data); err != nil {
		log.Printf("Write json file filename: %s , error: %s", s.file.Name(), err)
	}
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
