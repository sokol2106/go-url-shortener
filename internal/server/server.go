package server

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(handler http.Handler, addr string) *Server {
	return &Server{
		httpServer: &http.Server{
			Handler: handler,
			Addr:    addr,
		},
	}
}

func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
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
