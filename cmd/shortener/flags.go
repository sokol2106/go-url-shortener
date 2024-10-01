// Package main предоставляет функции для парсинга флагов командной строки.
package main

import "flag"

// ParseFlags разбирает флаги командной строки и обновляет параметры конфигурации.
// Принимает указатель на структуру params, которая содержит конфигурационные параметры.
// Флаги:
//
//	-a: адрес для запуска сервера (по умолчанию используется значение из params)
//	-b: базовый адрес для сокращенных URL (по умолчанию используется значение из params)
//	-f: путь к файлу хранения (по умолчанию используется значение из params)
//	-d: строка подключения к базе данных (по умолчанию используется значение из params)
func ParseFlags(p *params) {
	flag.StringVar(&p.ServerAddress, "a", p.ServerAddress, "address to run server")
	flag.StringVar(&p.BaseAddress, "b", p.BaseAddress, "base address of the resulting shortened URL")
	flag.StringVar(&p.FileStoragePath, "f", p.FileStoragePath, "file storage path")
	flag.StringVar(&p.DatabaseDSN, "d", p.DatabaseDSN, "data connection Database")
	flag.Parse()
}
