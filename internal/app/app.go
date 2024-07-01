package app

import (
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/server"
	"github.com/sokol2106/go-url-shortener/internal/storage"
	"log"
)

func Run(bsCnf *config.ConfigServer, shCnf *config.ConfigServer, fileStoragePath string) {
	var strg storage.ShortDataList
	strg.Init(fileStoragePath)
	sh := shorturl.New(shCnf.URL(), strg)
	ser := server.NewServer(shorturl.Router(sh), bsCnf.Addr())
	err := ser.Start()
	if err != nil {
		log.Printf("Starting server error: %s", err)
	}
}
