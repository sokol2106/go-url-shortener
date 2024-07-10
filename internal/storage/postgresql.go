package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"log"
	"time"
)

/*
type DataBase interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	PingContext(ctx context.Context) error
	Close()
}
*/

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
	if err != nil {
		log.Println("error connecting to Postgresql ", err)
		return err
	}

	return nil
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

	log.Printf("get Postgresql: %s", shURL)

	var originalURL string
	rows := pstg.db.QueryRowContext(context.Background(), "SELECT original FROM public.shorturl WHERE short=$1", shURL)
	err := rows.Scan(&originalURL)
	if err != nil {
		log.Println("error scanning short url postgresql", err)
		return ""
	}

	return originalURL
}

func (pstg *Postgresql) AddURL(originalURL string) string {
	var shortURL string
	hash := GenerateHash(originalURL)
	rows := pstg.db.QueryRowContext(context.Background(),
		"INSERT INTO public.shorturl (key, short, original) VALUES "+
			" ($1 ,$2 ,$3 ) ON CONFLICT (original) "+
			"DO UPDATE SET original = EXCLUDED.original RETURNING short",
		hash,
		RandText(8),
		originalURL)

	err := rows.Scan(&shortURL)
	if err != nil {
		log.Println("error adding short url postgresql", err)
		return ""
	}

	return shortURL
}

func (pstg *Postgresql) AddBatch(req []shorturl.RequestBatch, redirectURL string) []shorturl.ResponseBatch {
	resp := make([]shorturl.ResponseBatch, len(req))
	for i, val := range req {
		sh := pstg.AddURL(val.OriginalURL)
		log.Printf("ADD path: %s, URL: %s", sh, val.OriginalURL)
		resp[i] = shorturl.ResponseBatch{CorrelationID: val.CorrelationID, ShortURL: fmt.Sprintf("%s/%s", redirectURL, sh)}
	}
	return resp
}

func (pstg *Postgresql) Migrations(pathFiles string) error {
	driver, err := postgres.WithInstance(pstg.db, &postgres.Config{})
	if err != nil {
		log.Printf("error creating postgres driver: %v", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(pathFiles, "postgres", driver)
	if err != nil {
		log.Println("error migrate Postgresql", err)
		return err
	}

	if err = m.Up(); err != nil {
		log.Println("error up Postgresql", err)
		return err
	}

	return nil
}
