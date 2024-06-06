package main

import (
	"github.com/sokol2106/go-url-shortener/internal/app"
	"github.com/sokol2106/go-url-shortener/internal/config"
	"os"
)

type params struct {
	SERVER_ADDRESS string
	BASE_URL       string
}

func main() {
	var (
		flg  bool
		flg2 bool
	)

	p := params{"localhost:8080", "localhost:8080"}
	p.SERVER_ADDRESS, flg = os.LookupEnv("SERVER_ADDRESS")
	p.BASE_URL, flg2 = os.LookupEnv("BASE_URL")

	if !flg || !flg2 {
		ParseFlags(&p.SERVER_ADDRESS, &p.BASE_URL)
	}
	
	app.Run(config.NewConfigURL(p.SERVER_ADDRESS), config.NewConfigURL(p.BASE_URL))
}
