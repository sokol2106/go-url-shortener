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
const CBaseURL = "http://localhost"

// СFileStoragePath - путь к файлу хранения по умолчанию.
const СFileStoragePath = "/tmp/short-url-db.json"

// ConfigServer представляет конфигурацию сервера, включая хост и порт.
// "server_address": "localhost:8080", аналог переменной окружения SERVER_ADDRESS или флага -a
// "base_url": "http://localhost", аналог переменной окружения BASE_URL или флага -b
// "file_storage_path": "/path/to/file.db", аналог переменной окружения FILE_STORAGE_PATH или флага -f
// "database_dsn": "", аналог переменной окружения DATABASE_DSN или флага -d
// "enable_https": true аналог переменной окружения ENABLE_HTTPS или флага -s
type ConfigServer struct {
	serverAddress   string
	baseUrl         string
	fileStoragePath string
	databaseDsn     string
	enableHttps     bool
}

// NewConfigURL создает новый экземпляр ConfigServer на основе переданного URL.
// Возвращает указатель на ConfigServer и ошибку, если парсинг URL не удался.
func NewConfigURL(serverAddress, baseURL, fileStoragePath, databaseDsn, enableHttps string) *ConfigServer {
	enHttps, err := strconv.ParseBool(enableHttps)
	if err != nil {
		enHttps = false
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

	return &ConfigServer{
		serverAddress:   serverAddress,
		baseUrl:         baseURL,
		fileStoragePath: fileStoragePath,
		databaseDsn:     databaseDsn,
		enableHttps:     enHttps,
	}
}

// NewConfig создает новый экземпляр ConfigServer с пустыми полями.
func NewConfig() *ConfigServer {
	return &ConfigServer{}
}

// SetServerAddress задает адрес сервера
func (cs *ConfigServer) SetServerAddress(serverAddress string) *ConfigServer {
	cs.serverAddress = serverAddress
	return cs
}

// SetBaseUrl задает базовый URL.
func (cs *ConfigServer) SetBaseUrl(baseUrl string) *ConfigServer {
	cs.baseUrl = baseUrl
	return cs
}

// SetFileStoragePath задает путь к файлу хранилища.
func (cs *ConfigServer) SetFileStoragePath(fileStoragePath string) *ConfigServer {
	cs.fileStoragePath = fileStoragePath
	return cs
}

// SetDatabaseDsn задает строку подключения к базе данных.
func (cs *ConfigServer) SetDatabaseDsn(databaseDsn string) *ConfigServer {
	cs.databaseDsn = databaseDsn
	return cs
}

// SetEnableHttps задает флаг включения HTTPS, принимает строку и преобразует её в bool.
func (cs *ConfigServer) SetEnableHttps(enableHttps string) *ConfigServer {
	enHttps, err := strconv.ParseBool(enableHttps)
	if err != nil {
		enHttps = false
	}
	cs.enableHttps = enHttps
	return cs
}

// GetServerAddress возвращает адрес сервера в формате "host:port".
func (cs *ConfigServer) ServerAddress() string {
	return cs.serverAddress
}

// GetServerURL возвращает полный URL сервера
func (cs *ConfigServer) ServerURL() string {
	if cs.enableHttps {
		return fmt.Sprintf("https://%s", cs.serverAddress)
	}
	return fmt.Sprintf("https://%s", cs.serverAddress)
}

// EnableHTTPS возвращает флаг включения HTTPS.
func (cs *ConfigServer) EnableHTTPS() bool {
	return cs.enableHttps
}

// BaseUrl возвращает базовый URL
func (cs *ConfigServer) BaseUrl() string {
	return cs.baseUrl
}

// DatabaseDsn параметр подключения к БД
func (cs *ConfigServer) DatabaseDsn() string {
	return cs.databaseDsn
}

// FileStoragePath возвращает путь к файлу для сохранения
func (cs *ConfigServer) FileStoragePath() string {
	return cs.fileStoragePath
}
