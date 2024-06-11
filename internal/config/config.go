package config

import (
	"errors"
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

func CheckURL(body string) error {
	urlParse, err := url.Parse(body)
	if err != nil {
		return err
	}

	if urlParse.Scheme != "http" && urlParse.Scheme != "https" || urlParse.Host == "" {
		return errors.New("invalid url")
	}

	return nil
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
