// Package main предоставляет функции для парсинга флагов командной строки.
package main

import (
	"flag"
	"fmt"
)

// Option представляет собой функциональную опцию, которая используется для
// настройки конфигурации приложения через парсинг флагов командной строки.
type Option func()

// WithServerAddress создает опцию для задания параметров серверной конфигурации.
// Эти параметры включают адрес сервера, базовый адрес для сокращенных URL,
// путь к файлам хранилища и строку подключения к базе данных.
//
// Принимает указатель на структуру params, где будут сохраняться значения.
func WithServerAddress(p *params) Option {
	return func() {
		flag.StringVar(&p.ServerAddress, "a", p.ServerAddress, "address to run server")
		flag.StringVar(&p.BaseAddress, "b", p.BaseAddress, "base address of the resulting shortened URL")
		flag.StringVar(&p.FileStoragePath, "f", p.FileStoragePath, "file storage path")
		flag.StringVar(&p.DatabaseDSN, "d", p.DatabaseDSN, "data connection Database")
	}
}

// WithBuildInfo создает опцию для задания информации о сборке приложения.
// Эти параметры включают версию сборки, дату сборки и хэш коммита.
func WithBuildInfo() Option {
	return func() {
		flag.StringVar(&buildInfo.Version, "buildVersion", buildInfo.Version, "version of the build ")
		flag.StringVar(&buildInfo.Date, "buildDate", buildInfo.Date, "date when the build was created")
		flag.StringVar(&buildInfo.Commit, "buildCommit", buildInfo.Commit, "commit hash of the build")
	}
}

// ParseFlags разбирает флаги командной строки и обновляет параметры конфигурации
// на основе переданных функциональных опций. Каждая опция — это функция, которая
// настраивает определенные аспекты параметров конфигурации.
func ParseFlags(opts ...Option) {
	for _, opt := range opts {
		opt()
	}
	flag.Parse()
}

// printBuildInfo выводит на экран информацию о версии сборки, дате сборки
// и хэше коммита, взятых из структуры buildInfo.
func printBuildInfo() {
	fmt.Printf("Build version: %s\nBuild date: %s\nBuild commit: %s\n",
		buildInfo.Version, buildInfo.Date, buildInfo.Commit)
}
