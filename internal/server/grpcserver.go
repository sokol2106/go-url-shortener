package server

import (
	"github.com/sokol2106/go-url-shortener/internal/handlers/grpchandlers"
	"github.com/sokol2106/go-url-shortener/internal/handlers/grpchandlers/grpcservice/proto"
	"github.com/sokol2106/go-url-shortener/internal/service"
	"google.golang.org/grpc"
	"net"
)

// grpcServer структура gRPC-сервера.
type grpcServer struct {
	server *grpc.Server
}

type Option func(*grpcServer)

// WithMaxConnections создает gRPC-сервер с параметром максимальным количеством подключений.
func WithMaxConnections(max int) Option {
	return func(s *grpcServer) {
		s.server = grpc.NewServer(
			grpc.MaxConcurrentStreams(uint32(max)),
		)
	}
}

// NewGRPCServer создает gRPC-сервер.
func NewGRPCServer(opt ...Option) *grpcServer {
	s := &grpcServer{
		server: grpc.NewServer(),
	}

	for _, option := range opt {
		option(s)
	}

	return s
}

// StartGRPCServer создает и запускает gRPC-сервер.
func (g *grpcServer) StartGRPCServer(address string, srvSh *service.ShortURL, srvAu *service.Authorization, subnet string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	server := grpchandlers.NewURLShortenerServer(srvSh, srvAu, subnet)
	proto.RegisterURLShortenerServer(g.server, server)

	return g.server.Serve(lis)
}

// StopGRPCServer останавливает gRPC-сервер.
func (g *grpcServer) StopGRPCServer() {
	g.server.GracefulStop()
}
