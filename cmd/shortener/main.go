package main

import (
	"github.com/sokol2106/go-url-shortener/internal/app"
	"github.com/sokol2106/go-url-shortener/internal/config"
)

func main() {

	var (
		add  string
		add2 string
	)

	ParseFlags(&add, &add2)
	app.Run(config.NewConfigURL(add), config.NewConfigURL(add2))
}
