package shorturl

import "github.com/sokol2106/go-url-shortener/internal/storage"

// Для дальнейшей модификации
type StorageURL interface {
	AddURL(url string) string
	GetURL() string
	Close() error
}

type Database interface {
	Connect() error
	Disconnect() error
	PingContext() error
}

type ShortURL struct {
	redirectURL string
	storageURL  storage.ShortDataList
	database    Database
}

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}
