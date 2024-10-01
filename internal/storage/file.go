// Package storage предоставляет реализацию интерфейса хранилища URL-ов.
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

// File представляет структуру для работы с файлами, реализующую интерфейс Storage.
// file - это открытый файл для чтения и записи данных, scanner и decoder используются для чтения данных,
// encoder - для записи, data хранит считанные из файла данные.
type File struct {
	file          *os.File          // открытый файл для работы с данными
	scanner       *bufio.Scanner    // сканер для последовательного чтения файла
	decoder       *json.Decoder     // декодер JSON для чтения данных из файла
	encoder       *json.Encoder     // энкодер JSON для записи данных в файл
	data          []model.ShortData // кеш данных о сокращённых URL
	fileName      string            // имя файла для работы с URL
	isWriteEnable bool              // флаг, указывающий, разрешена ли запись в файл
}

// NewFile открывает файл для работы с URL-ами или создаёт новый файл, если он не существует.
// Возвращает указатель на созданную структуру File.
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

// AddOriginalURL добавляет оригинальный URL и возвращает сгенерированный сокращённый URL.
// Также сохраняет сокращённый URL в файл, если включена запись.
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

// AddOriginalURLBatch добавляет список URL-ов в пакетном режиме и возвращает соответствующие сокращённые URL-ы.
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

// GetOriginalURL ищет оригинальный URL по сокращённому. Если URL не найден, возвращает ошибку.
func (s *File) GetOriginalURL(ctx context.Context, shURL string) (model.ShortData, error) {
	var result model.ShortData
	ctxF, cancelF := context.WithCancel(ctx)
	defer cancelF()
	shortData := s.find(ctxF, model.ShortData{UUID: "", OriginalURL: "", ShortURL: shURL})
	if shortData != nil {
		return *shortData, nil
	}

	return result, cerrors.ErrGetShortURLNotFind

}

// GetUserShortenedURLs возвращает список сокращённых URL-ов пользователя, найденных в файле.
func (s *File) GetUserShortenedURLs(ctx context.Context, userID, redirectURL string) ([]service.ResponseUserShortenedURL, error) {
	result := make([]service.ResponseUserShortenedURL, 0)
	newFile, _ := os.OpenFile(s.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer newFile.Close()
	scanner := bufio.NewScanner(newFile)
	for scanner.Scan() {
		data := scanner.Bytes()
		mdl := &model.ShortData{}
		if err := json.Unmarshal(data, mdl); err != nil {
			return nil, err
		}

		if mdl.UserID == userID {
			result = append(result, service.ResponseUserShortenedURL{OriginalURL: mdl.OriginalURL, ShortURL: fmt.Sprintf("%s/%s", redirectURL, mdl.ShortURL)})
		}

	}

	return result, nil
}

// DeleteOriginalURL не реализован. Возвращает nil.
func (s *File) DeleteOriginalURL(ctx context.Context, data service.RequestUserShortenedURL) error {
	return nil
}

// getOrCreateShortData ищет данные о сокращённом URL или создаёт их, если они не найдены.
// Возвращает найденные или созданные данные, а также флаг, указывающий, были ли данные новыми.
func (s *File) getOrCreateShortData(ctx context.Context, hash, url, userID string) (*model.ShortData, bool) {
	ctxF, cancelF := context.WithCancel(ctx)
	defer cancelF()
	shortData := s.find(ctxF, model.ShortData{UUID: hash, OriginalURL: url, ShortURL: ""})
	isNewShortData := false
	if shortData == nil {
		isNewShortData = true
		shortData = &model.ShortData{UUID: hash, ShortURL: RandText(8), OriginalURL: url, UserID: userID, DeletedFlag: false}
	}

	return shortData, isNewShortData
}

// find ищет данные о сокращённом URL в файле на основе переданных параметров.
func (s *File) find(ctx context.Context, shortData model.ShortData) *model.ShortData {
	newFile, _ := os.OpenFile(s.fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	defer newFile.Close()
	scanner := bufio.NewScanner(newFile)
	for scanner.Scan() {
		data := scanner.Bytes()
		sd := &model.ShortData{}
		if err := json.Unmarshal(data, sd); err != nil {
			fmt.Printf("error Unmarshal file error: %s", err)
			return nil
		}

		if shortData.ShortURL == sd.ShortURL || shortData.UUID == sd.UUID || shortData.OriginalURL == sd.OriginalURL {
			return sd
		}
	}

	return nil
}

// PingContext не реализован. Возвращает ошибку.
func (s *File) PingContext() error {
	return errors.New("ping FILE not yet implemented ")
}

// Close закрывает файл.
func (s *File) Close() error {
	return s.file.Close()
}
