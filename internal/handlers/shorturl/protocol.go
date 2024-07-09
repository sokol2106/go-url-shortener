package shorturl

// Для дальнейшей модификации
type StorageURL interface {
	AddURL(string) string
	AddBatch([]RequestBatch) []ResponseBatch
	GetURL(string) string
	PingContext() error
	Close() error
}

type ShortURL struct {
	redirectURL string
	storageURL  StorageURL
}

type RequestJSON struct {
	URL string `json:"url"`
}

type ResponseJSON struct {
	Result string `json:"result"`
}

type RequestBatch struct {
	CorrelationId string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseBatch struct {
	CorrelationId string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
