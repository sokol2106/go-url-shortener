// Package config предоставляет структуру и функции для работы с конфигурацией сервера.
// Он включает в себя парсинг URL для извлечения хоста и порта, а также методы для получения
// информации о конфигурации сервера.
package config

import (
	"fmt"
	"strconv"
)

// CServerAddress - адрес сервера по умолчанию.
const CServerAddress = "localhost:8080"

// CBaseAddress - базовый адрес по умолчанию.
const CBaseURL = "http://localhost:8080"

// СFileStoragePath - путь к файлу хранения по умолчанию.
const СFileStoragePath = "/tmp/short-url-db.json"

// ConfigServer представляет конфигурацию сервера, включая хост и порт.
// "server_address": "localhost:8080", аналог переменной окружения SERVER_ADDRESS или флага -a
// "base_url": "http://localhost", аналог переменной окружения BASE_URL или флага -b
// "file_storage_path": "/path/to/file.db", аналог переменной окружения FILE_STORAGE_PATH или флага -f
// "database_dsn": "", аналог переменной окружения DATABASE_DSN или флага -d
// "enable_https": true аналог переменной окружения ENABLE_HTTPS или флага -s
type ConfigServer struct {
	ServerAddress   string
	BaseURL         string
	FileStoragePath string
	DatabaseDsn     string
	EnableHTTPS     bool
}

// NewConfigURL создает новый экземпляр ConfigServer на основе переданного URL.
// Возвращает указатель на ConfigServer и ошибку, если парсинг URL не удался.
func NewConfigURL(serverAddress, baseURL, fileStoragePath, databaseDsn, enable string) *ConfigServer {
	enHTTPS, err := strconv.ParseBool(enable)
	if err != nil {
		enHTTPS = false
	}

	if serverAddress == "" {
		serverAddress = CServerAddress
	}
	if baseURL == "" {
		baseURL = CBaseURL
	}
	if fileStoragePath == "" {
		fileStoragePath = СFileStoragePath
	}

	//	urlParse, err := url.Parse(baseURL)
	//if err != nil {
	//		return nil
	//	}

	return &ConfigServer{
		ServerAddress:   serverAddress,
		BaseURL:         baseURL, //fmt.Sprintf("http://%s:%s", urlParse.Scheme, urlParse.Opaque),
		FileStoragePath: fileStoragePath,
		DatabaseDsn:     databaseDsn,
		EnableHTTPS:     enHTTPS,
	}
}

// SetEnableHttps задает флаг включения HTTPS, принимает строку и преобразует её в bool.
func (cs *ConfigServer) SetEnableHTTPS(enable string) *ConfigServer {
	enHTTPS, err := strconv.ParseBool(enable)
	if err != nil {
		enHTTPS = false
	}
	cs.EnableHTTPS = enHTTPS
	return cs
}

// GetServerURL возвращает полный URL сервера
func (cs *ConfigServer) ServerURL() string {
	if cs.EnableHTTPS {
		return fmt.Sprintf("https://%s", cs.ServerAddress)
	}
	return fmt.Sprintf("https://%s", cs.ServerAddress)
}

// GetServerURL возвращает полный URL сервера
func (cs *ConfigServer) GetServerURL() string {
	if cs.EnableHTTPS {
		return fmt.Sprintf("https://%s", cs.ServerAddress)
	}
	return fmt.Sprintf("https://%s", cs.ServerAddress)
}
