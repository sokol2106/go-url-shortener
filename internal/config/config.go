// Package config предоставляет структуру и функции для работы с конфигурацией сервера.
// Он включает в себя парсинг URL для извлечения хоста и порта, а также методы для получения
// информации о конфигурации сервера.
package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
)

// CServerAddress - адрес сервера по умолчанию.
const CServerAddress = "localhost:8080"

// CBaseAddress - базовый адрес по умолчанию.
const DefaultBaseURL = "http://localhost:8080"

// СFileStoragePath - путь к файлу хранения по умолчанию.
const СFileStoragePath = "/tmp/short-url-db.json"

// ConfigServer представляет конфигурацию сервера, включая хост и порт.
// "server_address": "localhost:8080", аналог переменной окружения SERVER_ADDRESS или флага -a
// "base_url": "http://localhost", аналог переменной окружения BASE_URL или флага -b
// "file_storage_path": "/path/to/file.db", аналог переменной окружения FILE_STORAGE_PATH или флага -f
// "database_dsn": "", аналог переменной окружения DATABASE_DSN или флага -d
// "enable_https": true аналог переменной окружения ENABLE_HTTPS или флага -s
type ConfigServer struct {
	ServerAddress   string `json:"server_address"`
	DefaultBaseURL  string `json:"default_base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DatabaseDSN     string `json:"database_dsn"`
	EnableHTTPS     bool   `json:"enable_https"`
	TrustedSubnet   string `json:"trusted_subnet"`
}

// NewConfigURL создает новый экземпляр ConfigServer на основе переданного URL.
// Возвращает указатель на ConfigServer и ошибку, если парсинг URL не удался.
func NewConfigURL(serverAddress, defaultBaseURL, fileStoragePath, databaseDsn, enable, trustedSubnet string) *ConfigServer {
	enHTTPS, err := strconv.ParseBool(enable)
	if err != nil {
		enHTTPS = false
	}

	if serverAddress == "" {
		serverAddress = CServerAddress
	}
	if defaultBaseURL == "" {
		defaultBaseURL = DefaultBaseURL
	}
	if fileStoragePath == "" {
		fileStoragePath = СFileStoragePath
	}

	return &ConfigServer{
		ServerAddress:   serverAddress,
		DefaultBaseURL:  defaultBaseURL,
		FileStoragePath: fileStoragePath,
		DatabaseDSN:     databaseDsn,
		EnableHTTPS:     enHTTPS,
		TrustedSubnet:   trustedSubnet,
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

// LoadFileConfig загружает файл конфигурации и обрабатывает его.
func (cs *ConfigServer) LoadFileConfig(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Fatal(err)
	}

	var cnf ConfigServer

	err = json.Unmarshal(data, &cnf)
	if err != nil {
		return err
	}

	cs.ServerAddress = cnf.ServerAddress
	cs.DefaultBaseURL = cnf.DefaultBaseURL
	cs.FileStoragePath = cnf.FileStoragePath
	cs.DatabaseDSN = cnf.DatabaseDSN
	cs.EnableHTTPS = cnf.EnableHTTPS
	cs.TrustedSubnet = cnf.TrustedSubnet
	return nil
}
