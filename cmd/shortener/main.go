// Package main предоставляет основную точку входа для приложения URL Shortener.
// Настраивает и запускает HTTP-сервер, а также включает поддержку профилирования через pprof.
package main

import (
	"github.com/sokol2106/go-url-shortener/internal/app"
	"github.com/sokol2106/go-url-shortener/internal/config"
	"log"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof
	"os"
)

// CServerAddress - адрес сервера по умолчанию.
const CServerAddress = "localhost:8080"

// CBaseAddress - базовый адрес по умолчанию.
const CBaseAddress = "localhost:8080"

// СFileStoragePath - путь к файлу хранения по умолчанию.
const СFileStoragePath = "/tmp/short-url-db.json"

// params представляет параметры конфигурации для приложения.
type params struct {
	ServerAddress   string
	BaseAddress     string
	FileStoragePath string
	DatabaseDSN     string
	EnableHTTPS     string
}

// BuildInfo представляет информацию о приложении
type BuildInfo struct {
	Version string // buildVersion версия сборки
	Date    string // buildDate дата сборки
	Commit  string // buildCommit комментарий сборки
}

var buildInfo = BuildInfo{
	Version: "N/A",
	Date:    "N/A",
	Commit:  "N/A",
}

// main является основной точкой входа приложения.
func main() {

	// Запускаем pprof для профилирования на порту 6060
	go func() {
		http.ListenAndServe("localhost:6060", nil) // запускаем pprof на 6060 порту
	}()

	// Получаем параметры конфигурации из переменных окружения
	p := params{
		ServerAddress:   os.Getenv("SERVER_ADDRESS"),
		BaseAddress:     os.Getenv("BASE_URL"),
		FileStoragePath: os.Getenv("FILE_STORAGE_PATH"),
		DatabaseDSN:     os.Getenv("DATABASE_DSN"),
		EnableHTTPS:     os.Getenv("ENABLE_HTTPS"),
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
	ParseFlags(WithServerAddress(&p), WithBuildInfo())
	printBuildInfo()

	configServer, err := config.NewConfigURL(p.ServerAddress, p.EnableHTTPS)
	if err != nil {
		log.Printf("Creating server config error: %s", err.Error())
		return
	}
	configBase, err := config.NewConfigURL(p.BaseAddress, "")
	if err != nil {
		log.Printf("Creating server config base address error: %s", err.Error())
		return
	}

	// Запускаем приложение с заданными параметрами конфигурации
	app.Run(configServer, configBase, app.WithDatabase(p.DatabaseDSN), app.WithFile(p.FileStoragePath))
}
