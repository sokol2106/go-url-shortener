package app

import (
	"github.com/sokol2106/go-url-shortener/internal/handlers/shorturl"
	"github.com/sokol2106/go-url-shortener/internal/server"
)

func Run() {
	srv := server.NewServer(shorturl.ShortRouter())
	err := srv.Start()
	if err != nil {

		panic(err)
	}
}
