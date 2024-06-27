package shorturl

import (
	"github.com/sokol2106/go-url-shortener/internal/storage"
)

type ShortURL struct {
	redirectURL   string
	shortDataList storage.ShortDatalList
}

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}
