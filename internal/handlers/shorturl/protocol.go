package shorturl

// Для дальнейшей модификации
type StorageURL interface {
	AddURL(string) string
	GetURL(string) string
	Close() error
}

type Database interface {
	PingContext() error
	Close() error
}

type ShortURL struct {
	redirectURL string
	storageURL  StorageURL
	database    Database
}

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}
