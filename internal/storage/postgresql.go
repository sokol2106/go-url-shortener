package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/sokol2106/go-url-shortener/internal/cerrors"
	"github.com/sokol2106/go-url-shortener/internal/model"
	"github.com/sokol2106/go-url-shortener/internal/service"
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

type PostgreSQL struct {
	db     *sql.DB
	config string
}

func NewPostgresql(cnf string) *PostgreSQL {
	var pstg = PostgreSQL{}
	pstg.config = cnf
	return &pstg
}

func (pstg *PostgreSQL) Connect() error {
	var err error
	pstg.db, err = sql.Open("pgx", pstg.config)
	if err != nil {
		log.Println("error connecting to Postgresql ", err)
		return err
	}

	err = pstg.PingContext()
	if err != nil {
		log.Println("error pinging Postgresql ", err)
		return err
	}

	return nil
}

func (pstg *PostgreSQL) Migrations(pathFiles string) error {
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

func (pstg *PostgreSQL) AddOriginalURL(originalURL, userID string) (string, error) {
	var err error

	err = nil
	hash := GenerateHash(originalURL)
	shortURL := RandText(8)
	_, errInser := pstg.db.ExecContext(context.Background(), "INSERT INTO public.shorturl (key, short, original, userid, deleteflag) VALUES ($1, $2, $3, $4, $5)",
		hash,
		shortURL,
		originalURL,
		userID,
		false)

	if errInser != nil {
		rows, errSelect := pstg.db.QueryContext(context.Background(), "SELECT short FROM public.shorturl WHERE key=$1", hash)
		if errSelect != nil {
			return "", errSelect
		}

		errSelect = rows.Err()
		if errSelect != nil {
			return "", errSelect
		}

		if rows.Next() {
			errScan := rows.Scan(&shortURL)
			if errScan != nil {
				return "", errScan
			}
		}
		err = cerrors.ErrNewShortURL
	}

	return shortURL, err
}

func (pstg *PostgreSQL) AddOriginalURLBatch(req []service.RequestBatch, redirectURL, userID string) ([]service.ResponseBatch, error) {
	var err error
	err = nil
	resp := make([]service.ResponseBatch, len(req))
	for i, val := range req {
		sh, errAdd := pstg.AddOriginalURL(val.OriginalURL, userID)
		if errAdd != nil {
			if !errors.Is(errAdd, cerrors.ErrNewShortURL) {
				return nil, errAdd
			}
			err = errAdd
		}
		resp[i] = service.ResponseBatch{CorrelationID: val.CorrelationID, ShortURL: fmt.Sprintf("%s/%s", redirectURL, sh)}
	}
	return resp, err
}

func (pstg *PostgreSQL) GetOriginalURL(ctx context.Context, shURL string) (model.ShortData, error) {
	var (
		originalURL string
		deleteFlag  bool
		result      model.ShortData
	)

	ctxDB, cancelDB := context.WithCancel(ctx)
	defer cancelDB()
	rows := pstg.db.QueryRowContext(ctxDB, "SELECT original, deleteflag FROM public.shorturl WHERE short=$1", shURL)
	err := rows.Scan(&originalURL, &deleteFlag)
	if err != nil {
		if err == sql.ErrNoRows {
			return result, cerrors.ErrGetShortURLNotFind
		} else {
			return result, err
		}
	}

	return result, nil
}

func (pstg *PostgreSQL) GetUserShortenedURLs(ctx context.Context, userID, redirectURL string) ([]service.ResponseUserShortenedURL, error) {
	result := make([]service.ResponseUserShortenedURL, 0)
	var originalURL, shortURL string
	ctxDB, cancelDB := context.WithCancel(ctx)
	defer cancelDB()
	rows, err := pstg.db.QueryContext(ctxDB, "SELECT original, short FROM public.shorturl WHERE userid=$1", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&originalURL, &shortURL)
		if err != nil {
			return nil, err
		}
		result = append(result, service.ResponseUserShortenedURL{OriginalURL: originalURL, ShortURL: fmt.Sprintf("%s/%s", redirectURL, shortURL)})
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (pstg *PostgreSQL) DeleteOriginalURL(ctx context.Context, data service.RequestUserShortenedURL) error {
	_, errInser := pstg.db.ExecContext(context.Background(), "UPDATE public.shorturl SET deleteflag = true WHERE short=$1 AND userid=$2",
		data.ShortURL,
		data.UserID)
	return errInser
}

func (pstg *PostgreSQL) PingContext() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	return pstg.db.PingContext(ctx)
}

func (pstg *PostgreSQL) Close() error {
	if pstg.db != nil {
		return pstg.db.Close()
	}
	return nil
}
