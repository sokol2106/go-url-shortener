// Package server предоставляет возможности для инициализации, запуска и остановки HTTP-сервера.
package server

import (
	"context"
	"net/http"
)

// Server представляет структуру HTTP-сервера.
type Server struct {
	httpServer *http.Server
}

// NewServer создаёт сервер и возвращает аддрес объекта.
func NewServer(handler http.Handler, addr string) *Server {
	return &Server{
		httpServer: &http.Server{
			Handler: handler,
			Addr:    addr,
		},
	}
}

// Start запускает HTTP-сервер. Сервер начинает слушать входящие запросы.
func (s *Server) Start() error {
	return s.httpServer.ListenAndServe()
}

// Stop останавливает HTTP-сервер с возможностью плавного завершения работы.
// ctx - контекст, который может использоваться для задания тайм-аута на завершение работы сервера.
func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
