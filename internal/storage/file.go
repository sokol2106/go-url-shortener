package storage

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sokol2106/go-url-shortener/internal/cerrors"
	"github.com/sokol2106/go-url-shortener/internal/model"
	"github.com/sokol2106/go-url-shortener/internal/service"
	"os"
	"time"
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
		return &resFile
	}

	return nil
}

func (s *File) AddOriginalURL(originalURL, userID string) (string, error) {
	ctxF, cancelF := context.WithTimeout(context.Background(), 5000*time.Second)
	defer cancelF()
	err := cerrors.ErrNewShortURL
	hash := GenerateHash(originalURL)
	shortData, isNewShortData := s.getOrCreateShortData(ctxF, hash, originalURL, userID)
	if isNewShortData {
		err = nil
		if s.isWriteEnable {
			errEncode := s.encoder.Encode(&shortData)
			if errEncode != nil {
				err = errEncode
			}
		}
	}

	return shortData.ShortURL, err
}

func (s *File) AddOriginalURLBatch(req []service.RequestBatch, redirectURL string, userID string) ([]service.ResponseBatch, error) {
	var err error = nil
	resp := make([]service.ResponseBatch, len(req))
	for i, val := range req {
		sh, errAdd := s.AddOriginalURL(val.OriginalURL, userID)
		if errAdd != nil {
			if !errors.Is(errAdd, cerrors.ErrNewShortURL) {
				return nil, errAdd
			}
			err = errAdd
		}
		resp[i] = service.ResponseBatch{CorrelationID: val.CorrelationID, ShortURL: fmt.Sprintf("%s/%s", redirectURL, sh)}
	}

	return resp, err
}

func (s *File) GetOriginalURL(ctx context.Context, shURL string) string {
	ctxF, cancelF := context.WithCancel(ctx)
	defer cancelF()
	shortData := s.find(ctxF, model.ShortData{UUID: "", OriginalURL: "", ShortURL: shURL})
	if shortData == nil {
		return ""
	}
	return shortData.OriginalURL
}

func (s *File) getOrCreateShortData(ctx context.Context, hash, url, userID string) (*model.ShortData, bool) {
	ctxF, cancelF := context.WithCancel(ctx)
	defer cancelF()
	shortData := s.find(ctxF, model.ShortData{UUID: hash, OriginalURL: url, ShortURL: ""})
	isNewShortData := false
	if shortData == nil {
		isNewShortData = true
		shortData = &model.ShortData{UUID: hash, ShortURL: RandText(8), OriginalURL: url, UserID: userID}
	}

	return shortData, isNewShortData
}

func (s *File) find(ctx context.Context, shortData model.ShortData) *model.ShortData {
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
