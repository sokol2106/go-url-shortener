package fileStorage

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/sokol2106/go-url-shortener/internal/model"
	"github.com/sokol2106/go-url-shortener/internal/storage/memoryStorage"
	"log"
	"os"
)

type FileStorage struct {
	file          *os.File
	scanner       *bufio.Scanner
	encoder       *json.Encoder
	isWriteEnable bool
}

func New(filename string) *FileStorage {
	resFile := FileStorage{}
	resFile.isWriteEnable = false
	if filename != "" {
		newFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("error load file filename: %s , error: %s", filename, err)
			return nil
		}

		resFile.isWriteEnable = true
		resFile.encoder = json.NewEncoder(newFile)
		resFile.file = newFile
		resFile.scanner = bufio.NewScanner(newFile)
	}

	return &resFile
}

func (s *FileStorage) GetURL(shURL string) string {
	//scanner := bufio.NewScanner(s.file)
	for s.scanner.Scan() {
		data := s.scanner.Bytes()
		sd := model.ShortData{}
		err := json.Unmarshal(data, &sd)
		if err != nil {
			fmt.Printf("error Unmarshal file error: %s", err)
			return ""
		}

		if shURL == sd.ShortURL {
			return sd.OriginalURL
		}
	}

	return ""
}

func (s *FileStorage) AddURL(originalURL string) string {
	hash := memoryStorage.GenerateHash(originalURL)
	shortData, isNewShortData := s.getOrCreateShortData(hash, originalURL)
	if shortData == nil {
		log.Printf("Error Create Short Data")
		return ""
	}
	if isNewShortData && s.isWriteEnable {
		if err := s.encoder.Encode(&shortData); err != nil {
			log.Printf("Write json file filename: %s , error: %s", s.file.Name(), err)
		}
	}

	return shortData.ShortURL
}

func (s *FileStorage) getOrCreateShortData(hash, url string) (*model.ShortData, bool) {
	var shortData model.ShortData
	//scanner := bufio.NewScanner(s.file)
	for s.scanner.Scan() {
		data := s.scanner.Bytes()
		err := json.Unmarshal(data, &shortData)
		if err != nil {
			fmt.Printf("error Unmarshal file error: %s", err)
			return nil, false
		}

		if hash == shortData.UUID {
			return &shortData, false
		}
	}

	shortData = model.ShortData{UUID: hash, ShortURL: memoryStorage.RandText(8), OriginalURL: url}
	return &shortData, true
}

func (s *FileStorage) Close() error {
	return s.file.Close()
}
