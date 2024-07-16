package app

import (
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/server"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"log"
)

type App struct {
	Db   *storage.PostgreSQL
	File *storage.File
}

type Option func(*App)

func WithDatabase(dsn string) Option {
	return func(a *App) {
		a.Db = storage.NewPostgresql(dsn)
	}
}

func WithFile(filename string) Option {
	return func(a *App) {
		a.File = storage.NewFile(filename)
	}
}

func Run(bsCnf, shCnf *config.ConfigServer, opts ...Option) {
	var (
		handlerShort *shorturl.ShortURL
	)

	app := &App{}
	for _, opt := range opts {
		opt(app)
	}

	if err := app.Db.Connect(); err != nil {
		if app.File != nil {
			handlerShort = shorturl.New(shCnf.URL(), app.File)
		} else {
			handlerShort = shorturl.New(shCnf.URL(), storage.NewMemory())
		}
	} else {
		err = app.Db.Migrations("file://migrations/postgresql")
		if err != nil {
			log.Printf("error Migrations db: %s", err)
		}
		handlerShort = shorturl.New(shCnf.URL(), app.Db)
	}

	ser := server.NewServer(shorturl.Router(handlerShort), bsCnf.Addr())
	err := ser.Start()
	if err != nil {
		log.Printf("Starting server error: %s", err)
	}
}
