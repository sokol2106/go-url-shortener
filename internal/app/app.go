package app

import (
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/server"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"log"
)

func Run(bsCnf *config.ConfigServer, shCnf *config.ConfigServer, fileStoragePath string, databaseDSN string) {

	var (
		handlerShort *shorturl.ShortURL
	)

	if databaseDSN != "" {
		db := storage.NewPostgresql(databaseDSN)
		err := db.Connect()
		if err != nil {
			log.Printf("Error connect db: %s", err)
		}

		handlerShort = shorturl.New(shCnf.URL(), db)

	} else {
		if fileStoragePath != "" {
			objectStorage := storage.NewFile(fileStoragePath)
			handlerShort = shorturl.New(shCnf.URL(), objectStorage)
		} else {
			objectStorage := storage.NewMemory()
			handlerShort = shorturl.New(shCnf.URL(), objectStorage)
		}
	}

	ser := server.NewServer(shorturl.Router(handlerShort), bsCnf.Addr())
	err := ser.Start()
	if err != nil {
		log.Printf("Starting server error: %s", err)
	}
}
