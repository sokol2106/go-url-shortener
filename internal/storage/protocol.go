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

type ShortDataList struct {
	mapData       map[string]ShortData
	file          *os.File
	encoder       *json.Encoder
	isWriteEnable bool
}
