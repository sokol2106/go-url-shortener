package storage

import (
	"encoding/json"
	"os"
)

type ShortData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"shortURL"`
	OriginalURL string `json:"originalURL"`
}

type ShortDatalList struct {
	mapData       map[string]*ShortData
	file          *os.File
	encoder       *json.Encoder
	flagWriteFile bool
}
