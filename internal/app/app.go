package app

import (
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/sokol2106/go-url-shortener/internal/handlers"
	"github.com/sokol2106/go-url-shortener/internal/server"
	"github.com/sokol2106/go-url-shortener/internal/service"
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

func initStorage(db *storage.PostgreSQL, file *storage.File) service.StorageURL {
	if err := db.Connect(); err == nil {
		err = db.Migrations("file://migrations/postgresql")
		log.Printf("error Migrations db: %s", err)
		return db
	}

	if file != nil {
		return file
	} else {
		return storage.NewMemory()
	}

}

func Run(bsCnf, shCnf *config.ConfigServer, opts ...Option) {

	app := &App{}
	for _, opt := range opts {
		opt(app)
	}

	objStorage := initStorage(app.DB, app.File)
	srvShortURL := service.NewShortURL(shCnf.URL(), objStorage)
	handler := handlers.NewHandlers(srvShortURL)

	ser := server.NewServer(handlers.Router(handler), bsCnf.Addr())
	err := ser.Start()
	if err != nil {
		log.Printf("Starting server error: %s", err)
	}
}
