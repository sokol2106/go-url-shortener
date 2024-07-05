package storage

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type Postgresql struct {
	db     *sql.DB
	config string
}

func NewPostgresql(cnf string) *Postgresql {
	var pstg = Postgresql{}
	pstg.config = cnf

	return &pstg
}

func (pstg *Postgresql) Connect() error {
	var err error
	pstg.db, err = sql.Open("pgx", pstg.config)
	return err
}

func (pstg *Postgresql) Close() error {
	if pstg.db != nil {
		return pstg.db.Close()
	}
	return nil
}

func (pstg *Postgresql) PingContext() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return pstg.db.PingContext(ctx)
}

func (pstg *Postgresql) GetURL(shURL string) string {
	//for _, value := range s.mapData {
	//	if shURL == value.ShortURL {
	//		return value.OriginalURL
	//	}
	//	}
	return ""

}

func (pstg *Postgresql) AddURL(originalURL string) string {
	//	hash := GenerateHash(originalURL)
	//	shortData, isNewShortData := s.getOrCreateShortData(hash, originalURL)
	//	if isNewShortData && s.isWriteEnable {
	//		s.writeToFile(shortData)
	//	}

	//	return shortData.ShortURL
	return ""
}
