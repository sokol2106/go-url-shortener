package app

import (
	"fmt"
	"github.com/sokol2106/go-url-shortener/internal/config"
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/logger"
	"github.com/sokol2106/go-url-shortener/internal/server"
)

func Run(bsCnf *config.ConfigServer, shCnf *config.ConfigServer) {
	ser := server.NewServer(shorturl.ShortRouter(shCnf.URL()), bsCnf.Addr())
	err := ser.Start()
	if err != nil {
		fmt.Printf("error starting server: %s", err)
	}

	logger.Init()
}
