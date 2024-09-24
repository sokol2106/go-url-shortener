// Package config предоставляет структуру и функции для работы с конфигурацией сервера.
// Он включает в себя парсинг URL для извлечения хоста и порта, а также методы для получения
// информации о конфигурации сервера.
package config

import (
	"fmt"
	"net"
	"net/url"
)

// ConfigServer представляет конфигурацию сервера, включая хост и порт.
type ConfigServer struct {
	host string
	port string
}

// NewConfigURL создает новый экземпляр ConfigServer на основе переданного URL.
// Возвращает указатель на ConfigServer и ошибку, если парсинг URL не удался.
func NewConfigURL(u string) (*ConfigServer, error) {
	var (
		h string
		p string
	)
	urlParse, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("parsing URL error: url: %s , err: %s", u, err)
	}

	h = urlParse.Scheme
	p = urlParse.Opaque

	if (urlParse.Scheme == "http" || urlParse.Scheme == "https") && urlParse.Host != "" {
		h, p, _ = net.SplitHostPort(urlParse.Host)
	}

	return &ConfigServer{
		host: h,
		port: p,
	}, nil
}

// Host возвращает хост сервера.
func (cs *ConfigServer) Host() string {
	return cs.host
}

// Port возвращает порт сервера.
func (cs *ConfigServer) Port() string {
	return cs.port
}

// Addr возвращает адрес сервера в формате "host:port".
func (cs *ConfigServer) Addr() string {
	return fmt.Sprintf("%s:%s", cs.Host(), cs.Port())
}

// URL возвращает полный URL сервера в формате "http://host:port".
func (cs *ConfigServer) URL() string {
	return fmt.Sprintf("http://%s:%s", cs.Host(), cs.Port())
}
