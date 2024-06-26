package storage

import (
	"encoding/json"
	"os"
)

type ShortData struct {
	uuid         string
	short_url    string
	original_url string
}

type ShortDatalList struct {
	mapData       map[string]*ShortData
	file          *os.File
	encoder       *json.Encoder
	flagWriteFile bool
}
