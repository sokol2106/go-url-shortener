// Package main предоставляет основную точку входа для приложения URL Shortener.
// Настраивает и запускает HTTP-сервер, а также включает поддержку профилирования через pprof.
package main

import (
	"github.com/sokol2106/go-url-shortener/internal/app"
	"github.com/sokol2106/go-url-shortener/internal/config"
	"net/http"
	_ "net/http/pprof" // подключаем пакет pprof
	"os"
)

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

	cnf := config.NewConfigURL(
		os.Getenv("SERVER_ADDRESS"),
		os.Getenv("BASE_URL"),
		os.Getenv("FILE_STORAGE_PATH"),
		os.Getenv("DATABASE_DSN"),
		os.Getenv("ENABLE_HTTPS"),
	)

	ParseFlags(WithServerAddress(cnf), WithBuildInfo())
	printBuildInfo()

	// Запускаем приложение с заданными параметрами конфигурации
	app.Run(cnf, app.WithDatabase(cnf.DatabaseDsn()), app.WithFile(cnf.FileStoragePath()))
}
