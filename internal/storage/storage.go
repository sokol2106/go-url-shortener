package storage

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// ShortData

func NewShortData(uuid string, short_url string, original_url string) *ShortData {
	return &ShortData{uuid, short_url, original_url}
}

// ShortDatalList

func (s *ShortDatalList) Init(filename string) {
	s.mapData = make(map[string]*ShortData)
	s.flagWriteFile = false
	if filename != "" {
		newFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("error load file filename: %s , error: %s", filename, err)
			return
		}

		s.flagWriteFile = true
		s.encoder = json.NewEncoder(newFile)
		s.file = newFile
		s.LoadDateFile()
	}
}

func (s *ShortDatalList) Close() error {
	return s.file.Close()
}

func (s *ShortDatalList) LoadDateFile() {
	if !s.flagWriteFile {
		return
	}

	sd := &ShortData{}
	scanner := bufio.NewScanner(s.file)
	for scanner.Scan() {
		data := scanner.Bytes()
		err := json.Unmarshal(data, &sd)
		if err != nil {
			fmt.Printf("error Unmarshal file error: %s", err)
			return
		}
		s.mapData[sd.UUID] = sd
		sd = &ShortData{}
	}
}

func (s *ShortDatalList) AddURL(originalURL string) string {
	hash := sha256.Sum256([]byte(originalURL))
	thash := hex.EncodeToString(hash[:])
	tshdata, exist := s.mapData[thash]
	if !exist {
		tshdata = &ShortData{thash, RandText(8), originalURL}
		s.mapData[thash] = tshdata

		if s.flagWriteFile {
			if err := s.encoder.Encode(&tshdata); err != nil {
				fmt.Printf("error write json file filename: %s , error: %s", s.file.Name(), err)
			}
		}
	}

	return tshdata.ShortURL
}

func (s *ShortDatalList) GetURL(shURL string) string {
	for _, value := range s.mapData {
		if shURL == value.ShortURL {
			return value.OriginalURL
		}
	}
	return ""
}

func RandText(size int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]rune, size)
	for i := range b {
		b[i] = letterRunes[r.Intn(len(letterRunes))]
	}
	return string(b)
}
