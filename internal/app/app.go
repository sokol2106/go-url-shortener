package app

import (
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/server"
	"github.com/sokol2106/go-url-shortener/internal/storage/fileStorage"
	"github.com/sokol2106/go-url-shortener/internal/storage/memoryStorage"
	"github.com/sokol2106/go-url-shortener/internal/storage/postgresql"
	"log"
)

func Run(bsCnf *config.ConfigServer, shCnf *config.ConfigServer, fileStoragePath string, databaseDSN string) {

	var (
		handlerShort *shorturl.ShortURL
	)

	if databaseDSN != "" {
		db := postgresql.New(databaseDSN)
		err := db.Connect()
		if err != nil {
			log.Printf("Error connect db: %s", err)
		}

		handlerShort = shorturl.New(shCnf.URL(), db, db)

	} else {
		if fileStoragePath != "" {
			objectStorage := fileStorage.New(fileStoragePath)
			handlerShort = shorturl.New(shCnf.URL(), objectStorage, nil)
		} else {
			objectStorage := memoryStorage.New()
			handlerShort = shorturl.New(shCnf.URL(), objectStorage, nil)
		}
	}

	ser := server.NewServer(shorturl.Router(handlerShort), bsCnf.Addr())
	err := ser.Start()
	if err != nil {
		log.Printf("Starting server error: %s", err)
	}
}
