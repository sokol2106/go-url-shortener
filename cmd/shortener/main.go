package main

import (
	"fmt"
	"github.com/sokol2106/go-url-shortener/internal/app"
	"github.com/sokol2106/go-url-shortener/internal/config"
	"os"
)

const CServerAddress = "localhost:8080"
const CBaseAddress = "localhost:8080"

type params struct {
	ServerAddress string
	BaseAddress   string
}

func main() {
	p := params{
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
		BaseAddress:   os.Getenv("BASE_URL"),
	}
	if p.ServerAddress == "" {
		p.ServerAddress = CServerAddress
	}
	if p.BaseAddress == "" {
		p.BaseAddress = CBaseAddress
	}
	ParseFlags(&p.ServerAddress, &p.BaseAddress)
	configServer, err := config.NewConfigURL(p.ServerAddress)
	if err != nil {
		fmt.Println("error creating server config: %s", err.Error())
		return
	}
	configBase, err := config.NewConfigURL(p.BaseAddress)
	if err != nil {
		fmt.Println("error creating server config base address: %s", err.Error())
		return
	}
	app.Run(configServer, configBase)
}
