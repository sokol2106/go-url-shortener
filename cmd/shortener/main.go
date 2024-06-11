package main

import (
	"github.com/sokol2106/go-url-shortener/internal/app"
	"github.com/sokol2106/go-url-shortener/internal/config"
	"os"
)

type params struct {
	ServerAddress string
	BaseURL       string
}

func main() {
	var (
		flg  bool
		flg2 bool
	)

	p := params{"localhost:8080", "localhost:8080"}
	p.ServerAddress, flg = os.LookupEnv("SERVER_ADDRESS")
	p.BaseURL, flg2 = os.LookupEnv("BASE_URL")

	if !flg || !flg2 {
		ParseFlags(&p.ServerAddress, &p.BaseURL)
	}

	app.Run(config.NewConfigURL(p.ServerAddress), config.NewConfigURL(p.BaseURL))
}
