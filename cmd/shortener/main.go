package main

import (
	"github.com/sokol2106/go-url-shortener/internal/app"
	"github.com/sokol2106/go-url-shortener/internal/config"
	"os"
)

type params struct {
	Server_Address string
	Base_URL       string
}

func main() {
	var (
		flg  bool
		flg2 bool
	)

	p := params{"localhost:8080", "localhost:8080"}
	p.Server_Address, flg = os.LookupEnv("SERVER_ADDRESS")
	p.Base_URL, flg2 = os.LookupEnv("BASE_URL")

	if !flg || !flg2 {
		ParseFlags(&p.Server_Address, &p.Base_URL)
	}

	app.Run(config.NewConfigURL(p.Server_Address), config.NewConfigURL(p.Base_URL))
}
