package shorturl

import "github.com/sokol2106/go-url-shortener/internal/storage"

/*
type storageURL interface {
	AddURL(url string) string
	GetURL() string
	Close() error
}
*/

type ShortURL struct {
	redirectURL   string
	shortDataList storage.ShortDataList
}

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}
