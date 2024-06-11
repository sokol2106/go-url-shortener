package config

import (
	"fmt"
	"net"
	"net/url"
)

type ConfigServer struct {
	host string
	port string
}

func NewConfig(h string, p string) *ConfigServer {
	return &ConfigServer{
		host: h,
		port: p,
	}
}

func NewConfigURL(u string) *ConfigServer {
	var (
		h string
		p string
	)
	urlParse, _ := url.Parse(u)
	if urlParse.Scheme == "http" || urlParse.Scheme == "https" {
		h, p, _ = net.SplitHostPort(urlParse.Host)
	} else {
		// Если пришло localhost:8080
		h = urlParse.Scheme
		p = urlParse.Opaque
	}

	return &ConfigServer{
		host: h,
		port: p,
	}
}

func (cs *ConfigServer) Host() string {
	return cs.host
}

func (cs *ConfigServer) Port() string {
	return cs.port
}

func (cs *ConfigServer) Addr() string {
	return fmt.Sprintf("%s:%s", cs.Host(), cs.Port())
}

func (cs *ConfigServer) URL() string {
	return fmt.Sprintf("http://%s:%s", cs.Host(), cs.Port())
}

func (cs *ConfigServer) URLS() string {
	return fmt.Sprintf("https://%s:%s", cs.Host(), cs.Port())
}
