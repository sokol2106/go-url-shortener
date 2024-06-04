package config

import "fmt"

const (
	DefaultHTTPpPort = 8080
	DefaultHTTPHost  = "localhost"
)

type configServer struct {
	httpHost string
	httpPort int
}

func NewConfig(h string, p int) *configServer {
	return &configServer{
		httpHost: h,
		httpPort: p,
	}
}

func (cs *configServer) HTTPHost() string {
	return cs.httpHost
}

func (cs *configServer) HTTPPort() int {
	return cs.httpPort
}
func (cs *configServer) URL() string {
	return fmt.Sprintf("http://%s:%d", cs.HTTPHost(), cs.HTTPPort())
}
