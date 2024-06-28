package storage

import (
	"encoding/json"
	"github.com/sokol2106/go-url-shortener/internal/model"
	"os"
)

type ShortDataList struct {
	mapData       map[string]model.ShortData
	file          *os.File
	encoder       *json.Encoder
	isWriteEnable bool
}
