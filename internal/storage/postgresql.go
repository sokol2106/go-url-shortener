package storage

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/jackc/pgx/v5/stdlib"
	"log"
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
	if err != nil {
		return err
	}

	m, err := migrate.New(
		"file:https://github.com/sokol2106/go-url-shortener/blob/iter11/migrations/postgresql/tableshorturl.sql",
		pstg.config)

	if err != nil {
		log.Println("error migrate Postgresql", err)
	}

	if err = m.Up(); err != nil {
		log.Println("error up Postgresql", err)
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
	var originalURL string
	rows := pstg.db.QueryRowContext(context.Background(), "SELECT key, short, original FROM public.shorturl WHERE short=$1", shURL)
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
			" ($1,$2,$3) ON CONFLICT (original) "+
			"DO UPDATE SET original = EXCLUDED.original RETURNING short", hash, RandText(8), originalURL)

	err := rows.Scan(&shortURL)
	if err != nil {
		log.Println("error adding short url postgresql", err)
		return ""
	}

	return shortURL
}
