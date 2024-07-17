//go:generate mockgen -destination=../../../mocks/mock_shorturl.go -package=shorturl . StorageURL
package shorturl

import "context"

// Для дальнейшей модификации
type StorageURL interface {
	AddURL(string) (string, error)
	AddBatch([]RequestBatch, string) ([]ResponseBatch, error)
	GetURL(context.Context, string) string
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
	CorrelationID string `json:"correlation_id"`
	OriginalURL   string `json:"original_url"`
}

type ResponseBatch struct {
	CorrelationID string `json:"correlation_id"`
	ShortURL      string `json:"short_url"`
}
