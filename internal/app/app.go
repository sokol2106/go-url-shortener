package app

import (
	"fmt"
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/server"
)

func Run(bsCnf *config.ConfigServer, shCnf *config.ConfigServer, fileStoragePath string) {
	sh := shorturl.NewShortURL(shCnf.URL(), fileStoragePath)
	ser := server.NewServer(shorturl.ShortRouter(sh), bsCnf.Addr())
	err := ser.Start()
	if err != nil {
		fmt.Printf("error starting server: %s", err)
	}
}
