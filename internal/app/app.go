package app

import (
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/server"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"log"
)

type App struct {
	DB   *storage.PostgreSQL
	File *storage.File
}

type Option func(*App)

func WithDatabase(dsn string) Option {
	return func(a *App) {
		a.DB = storage.NewPostgresql(dsn)
	}
}

func WithFile(filename string) Option {
	return func(a *App) {
		a.File = storage.NewFile(filename)
	}
}

func initStorage(db *storage.PostgreSQL, file *storage.File) shorturl.StorageURL {
	if err := db.Connect(); err == nil {
		err = db.Migrations("file://migrations/postgresql")
		if err == nil {
			return db
		}
		log.Printf("error Migrations db: %s", err)
	}

	if file != nil {
		return file
	} else {
		return storage.NewMemory()
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

	objStorage := initStorage(app.DB, app.File)
	handlerShort = shorturl.New(shCnf.URL(), objStorage)

	/*if err := app.DB.Connect(); err != nil {
		if app.File != nil {
			handlerShort = shorturl.New(shCnf.URL(), app.File)
		} else {
			mem := storage.NewMemory()
			handlerShort = shorturl.New(shCnf.URL(), mem)
		}
	} else {
		err = app.DB.Migrations("file://migrations/postgresql")
		if err != nil {
			log.Printf("error Migrations db: %s", err)
		}
		handlerShort = shorturl.New(shCnf.URL(), app.DB)
	}*/

	ser := server.NewServer(shorturl.Router(handlerShort), bsCnf.Addr())
	err := ser.Start()
	if err != nil {
		log.Printf("Starting server error: %s", err)
	}
}
