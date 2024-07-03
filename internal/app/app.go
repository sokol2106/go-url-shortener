package app

import (
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/sokol2106/go-url-shortener/internal/database/postgresql"
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/server"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"log"
)

func Run(bsCnf *config.ConfigServer, shCnf *config.ConfigServer, fileStoragePath string, database_DSN string) {
	var (
		strg storage.ShortDataList
	)
	db := postgresql.New(database_DSN)
	err := db.Connect()
	if err != nil {
		log.Printf("Error connect db: %s", err)
	}

	strg.Init(fileStoragePath)
	sh := shorturl.New(shCnf.URL(), strg, db)
	ser := server.NewServer(shorturl.Router(sh), bsCnf.Addr())
	err = ser.Start()
	if err != nil {
		log.Printf("Starting server error: %s", err)
	}
}
