// Package model предоставляет реализации структур данных, используемых в системе сокращения URL.
// Эти структуры представляют модели данных, используемые для хранения и работы с URL-адресами.
package model

// ShortData представляет структуру данных для хранения информации о сокращенном URL.
// Поля структуры:
// - UUID: Уникальный идентификатор для записи.
// - ShortURL: Сгенерированный короткий URL.
// - OriginalURL: Оригинальный URL, который был сокращен.
// - UserID: Идентификатор пользователя, которому принадлежит сокращенный URL.
// - DeletedFlag: Флаг, указывающий, удален ли URL.
type ShortData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"shortURL"`
	OriginalURL string `json:"originalURL"`
	UserID      string `json:"userId"`
	DeletedFlag bool   `json:"is_deleted"`
}
