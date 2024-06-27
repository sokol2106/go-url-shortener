package storage

import (
	"encoding/json"
	"os"
)

type ShortData struct {
	Uuid         string `json:"uuid"`
	Short_url    string `json:"short_url"`
	Original_url string `json:"original_url"`
}

type ShortDatalList struct {
	mapData       map[string]*ShortData
	file          *os.File
	encoder       *json.Encoder
	flagWriteFile bool
}
