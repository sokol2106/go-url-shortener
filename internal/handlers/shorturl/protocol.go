package shorturl

import "github.com/sokol2106/go-url-shortener/internal/storage"

// Для дальнейшей модификации

type ShortURL struct {
	redirectURL string
	storageURL  storage.ShortDataList
}

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}
