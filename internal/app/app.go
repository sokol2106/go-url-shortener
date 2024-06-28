package app

import (
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/server"
	"log"
)

func Run(bsCnf *config.ConfigServer, shCnf *config.ConfigServer, fileStoragePath string) {
	sh := shorturl.NewShortURL(shCnf.URL(), fileStoragePath)
	ser := server.NewServer(shorturl.ShortRouter(sh), bsCnf.Addr())
	err := ser.Start()
	if err != nil {
		log.Printf("Starting server error: %s", err)
	}
}
