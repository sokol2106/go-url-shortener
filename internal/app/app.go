// Package app предоставляет основные функции для инициализации и запуска приложения.
package app

import (
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/sokol2106/go-url-shortener/internal/handlers"
	"github.com/sokol2106/go-url-shortener/internal/server"
	"github.com/sokol2106/go-url-shortener/internal/service"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"log"
)

// App представляет собой основную структуру приложения, содержащую компоненты
// для работы с хранилищем данных и файловой системой.
type App struct {
	DB   *storage.PostgreSQL
	File *storage.File
}

// Option представляет собой функцию, которая настраивает приложение.
type Option func(*App)

// WithDatabase создает опцию для установки PostgreSQL как хранилище данных приложения.
func WithDatabase(dsn string) Option {
	return func(a *App) {
		a.DB = storage.NewPostgresql(dsn)
	}
}

// WithFile создает опцию для установки файлового хранилища.
func WithFile(filename string) Option {
	return func(a *App) {
		a.File = storage.NewFile(filename)
	}
}

// initStorage инициализирует хранилище данных, используя PostgreSQL или файл.
func initStorage(db *storage.PostgreSQL, file *storage.File) service.Storage {
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

// Run инициализирует приложение и запускает HTTP-сервер.
// Принимает конфигурации для базового и сокращенного URL и опции для настройки хранилища данных.
func Run(bsCnf, shCnf *config.ConfigServer, opts ...Option) {

	app := &App{}
	for _, opt := range opts {
		opt(app)
	}

	objStorage := initStorage(app.DB, app.File)
	srvShortURL := service.NewShortURL(shCnf.URL(), objStorage)
	handler := handlers.NewHandlers(srvShortURL)

	ser := server.NewServer(handlers.Router(handler), bsCnf.Addr())
	err := ser.Start(bsCnf.EnableHTTPS())
	if err != nil {
		log.Printf("Starting server error: %s", err)
	}
}
