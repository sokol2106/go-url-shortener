package main

import (
	"github.com/sokol2106/go-url-shortener/internal/app"
	"github.com/sokol2106/go-url-shortener/internal/config"
	"log"
	"os"
)

const CServerAddress = "localhost:8080"
const CBaseAddress = "localhost:8080"
const СFileStoragePath = "/tmp/short-url-db.json"

type params struct {
	ServerAddress   string
	BaseAddress     string
	FileStoragePath string
	Database_DSN    string
}

func main() {
	p := params{
		ServerAddress:   os.Getenv("SERVER_ADDRESS"),
		BaseAddress:     os.Getenv("BASE_URL"),
		FileStoragePath: os.Getenv("FILE_STORAGE_PATH"),
		Database_DSN:    os.Getenv("FILE_STORAGE_PATH"),
	}
	if p.ServerAddress == "" {
		p.ServerAddress = CServerAddress
	}
	if p.BaseAddress == "" {
		p.BaseAddress = CBaseAddress
	}
	if p.FileStoragePath == "" {
		p.FileStoragePath = СFileStoragePath
	}
	ParseFlags(&p)
	configServer, err := config.NewConfigURL(p.ServerAddress)
	if err != nil {
		log.Printf("Creating server config error: %s", err.Error())
		return
	}
	configBase, err := config.NewConfigURL(p.BaseAddress)
	if err != nil {
		log.Printf("Creating server config base address error: %s", err.Error())
		return
	}
	app.Run(configServer, configBase, p.FileStoragePath, p.Database_DSN)
}
