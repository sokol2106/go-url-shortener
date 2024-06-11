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

func NewConfigURL(u string) (*ConfigServer, error) {
	var (
		h string
		p string
	)
	urlParse, err := url.Parse(u)
	if err != nil {
		return nil, fmt.Errorf("parsing URL error: url: %s , err: %s", u, err)
	}
	if urlParse.Scheme == "" {
		h = urlParse.Scheme
		p = urlParse.Opaque
	} else {
		if (urlParse.Scheme == "http" || urlParse.Scheme == "https") && urlParse.Host != "" {
			h, p, _ = net.SplitHostPort(urlParse.Host)
		} else {
			return nil, fmt.Errorf("protocol error: %s", urlParse.Scheme)
		}
	}
	return &ConfigServer{
		host: h,
		port: p,
	}, nil
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
