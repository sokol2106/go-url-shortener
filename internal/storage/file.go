package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/model"
	"log"
	"os"
)

type File struct {
	file          *os.File
	scanner       *bufio.Scanner
	decoder       *json.Decoder
	encoder       *json.Encoder
	data          []model.ShortData
	fileName      string
	isWriteEnable bool
}

func NewFile(filename string) *File {
	resFile := File{}
	resFile.isWriteEnable = false
	if filename != "" {
		newFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("error load file filename: %s , error: %s", filename, err)
			return nil
		}

		resFile.fileName = filename
		resFile.isWriteEnable = true
		resFile.encoder = json.NewEncoder(newFile)
		resFile.decoder = json.NewDecoder(newFile)
		resFile.file = newFile
		resFile.scanner = bufio.NewScanner(newFile)
	}

	return &resFile
}

func (s *File) AddURL(originalURL string) string {
	hash := GenerateHash(originalURL)
	shortData, isNewShortData := s.getOrCreateShortData(hash, originalURL)
	if isNewShortData && s.isWriteEnable {
		if err := s.encoder.Encode(&shortData); err != nil {
			log.Printf("Write json file filename: %s , error: %s", s.file.Name(), err)
		}
	}

	return shortData.ShortURL
}

func (s *File) AddBatch(req []shorturl.RequestBatch, redirectURL string) []shorturl.ResponseBatch {
	resp := make([]shorturl.ResponseBatch, len(req))
	for i, val := range req {
		sh := s.AddURL(val.OriginalURL)
		resp[i] = shorturl.ResponseBatch{CorrelationID: val.CorrelationID, ShortURL: fmt.Sprintf("%s/%s", redirectURL, sh)}
	}

	return resp
}

func (s *File) GetURL(shURL string) string {

	log.Printf("get file: %s", shURL)

	shortData := s.find(model.ShortData{UUID: "", OriginalURL: "", ShortURL: shURL})
	if shortData == nil {
		return ""
	}
	return shortData.OriginalURL
}

func (s *File) getOrCreateShortData(hash, url string) (*model.ShortData, bool) {
	shortData := s.find(model.ShortData{UUID: hash, OriginalURL: url, ShortURL: ""})
	isNewShortData := false
	if shortData == nil {
		isNewShortData = true
		shortData = &model.ShortData{UUID: hash, ShortURL: RandText(8), OriginalURL: url}
	}

	return shortData, isNewShortData
}

func (s *File) find(shortData model.ShortData) *model.ShortData {
	newFile, _ := os.OpenFile(s.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	scanner := bufio.NewScanner(newFile)
	for scanner.Scan() {
		data := scanner.Bytes()
		sd := &model.ShortData{}
		if err := json.Unmarshal(data, sd); err != nil {
			fmt.Printf("error Unmarshal file error: %s", err)
			newFile.Close()
			return nil
		}

		if shortData.ShortURL == sd.ShortURL || shortData.UUID == sd.UUID || shortData.OriginalURL == sd.OriginalURL {
			newFile.Close()
			return sd
		}
	}

	newFile.Close()
	return nil
}

func (s *File) PingContext() error {
	return errors.New("ping FILE not yet implemented ")
}

func (s *File) Close() error {
	return s.file.Close()
}
