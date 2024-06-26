package storage

import (
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

func (sd *ShortData) UUID() string {
	return sd.uuid
}

func (sd *ShortData) Short() string {
	return sd.short_url
}

func (sd *ShortData) Original() string {
	return sd.original_url
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
	}
}

func (s *ShortDatalList) Close() error {
	return s.file.Close()
}

func (s *ShortDatalList) LoadDateFile() {
	//reader := bufio.NewReader(s.file)
	//data, err := reader.ReadBytes('\n')
	//if err != nil {
	//		return
	//	}
}

func (s *ShortDatalList) AddURL(originalURL string) string {
	hash := sha256.Sum256([]byte(originalURL))
	thash := hex.EncodeToString(hash[:])
	tshdata, exist := s.mapData[thash]
	if !exist {
		tshdata = NewShortData(thash, RandText(8), originalURL)
		s.mapData[thash] = tshdata

		if s.flagWriteFile {
			//	data, err := json.Marshal(tshdata)
			//	if err != nil {
			//		fmt.Printf("error Marshal json , error: %s", err)
			//	}

			/*	if _, err := s.writerBuffer.Write(data); err != nil {
					fmt.Printf("error write json file filename: %s , error: %s", s.file.Name(), err)
				}

				if err := s.writerBuffer.WriteByte('\n'); err != nil {
					fmt.Printf("error write 'n' file filename: %s , error: %s", s.file.Name(), err)
				}

				s.writerBuffer.Flush()
			*/
		}
	}

	return tshdata.Short()
}

func (s *ShortDatalList) GetURL(shURL string) string {
	for _, value := range s.mapData {
		if shURL == value.Short() {
			return value.Original()
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
