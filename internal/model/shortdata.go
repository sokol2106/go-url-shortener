package model

type ShortData struct {
	UUID        string `json:"uuid"`
	ShortURL    string `json:"shortURL"`
	OriginalURL string `json:"originalURL"`
	UserID      string `json:"userId"`
}
